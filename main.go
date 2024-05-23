package main

//https://dynobase.dev/dynamodb-golang-query-examples/#query-with-sorting
import (
	"context"
	"cryptoGridProjectDemo/activityLogApp"
	"cryptoGridProjectDemo/binanceLib"
	"cryptoGridProjectDemo/internal/pkg/log"
	"cryptoGridProjectDemo/internal/v1/router"
	"cryptoGridProjectDemo/leverageApp"
	"cryptoGridProjectDemo/mainApp"
	"cryptoGridProjectDemo/orderApp"
	"cryptoGridProjectDemo/quoteApp"
	"cryptoGridProjectDemo/redisLib"
	"cryptoGridProjectDemo/wsLib"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"net/http"
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
}

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = false
)

var wsServe = func(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	c, _, err := websocket.DefaultDialer.Dial(cfg.Endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	c.SetReadLimit(655350)
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		if WebsocketKeepalive {
			keepAlive(c, WebsocketTimeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			handler(message)
		}
	}()
	return
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}

type WsQuoteRes struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	TradeTime int64  `json:"tradeTime"`
}

func backgroundTaskQuote(ctx context.Context, hub *wsLib.Hub, rdb *redis.Client) {
	fmt.Println("Start backgroundTaskQuote...")
	redisService := redisLib.NewRedisService(rdb)
	err := redisService.SetTaskHearthBeat("TaskQuote")
	if err != nil {
		fmt.Println(err)
	}

	wsAggTradeHandler := func(event *futures.WsAggTradeEvent) {

		channelName := "ws/quote/"
		wsQuoteRes := WsQuoteRes{Symbol: event.Symbol, Price: event.Price, TradeTime: event.TradeTime}
		value, err := json.Marshal(wsQuoteRes)
		if err != nil {

		}

		err2 := redisService.SetSymbolPrice(event.Symbol, event.Price)
		if err2 != nil {
			fmt.Println(err2)
		}

		hub.BroadcastChannel <- wsLib.BroadcastChannel{Channel: channelName, Message: value}
		//hub.Broadcast <- value

	}
	fmt.Println(wsAggTradeHandler)
	errHandler := func(err error) {
		fmt.Println(err)
	}
	//todo: 改成活的
	symbolList := map[string]string{
		"ETHUSDT":  "1s",
		"BTCUSDT":  "1s",
		"AVAXUSDT": "1s",
		"XLMUSDT":  "1s",
	}

	//symbols := []string{"BTCUSDT", "ETHUSDT", "AVAXUSDT"}
	//doneC, stopC, err := futures.WsCombinedAggTradeServe(symbols, wsAggTradeHandler, errHandler)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	wsAllMarkPriceHandler := func(event futures.WsAllMarkPriceEvent) {
		//fmt.Println(len(event))

		for _, e := range event {
			symbol := e.Symbol

			if _, ok := symbolList[symbol]; ok {
				price := e.MarkPrice
				//fmt.Println(event)

				channelName := "ws/quote/"
				wsQuoteRes := WsQuoteRes{Symbol: symbol, Price: price, TradeTime: e.Time}
				value, err := json.Marshal(wsQuoteRes)
				if err != nil {

				}

				err2 := redisService.SetSymbolPrice(symbol, price)
				if err2 != nil {
					fmt.Println(err2)
				}

				hub.BroadcastChannel <- wsLib.BroadcastChannel{Channel: channelName, Message: value}
				//hub.Broadcast <- value
			}
		}
	}
	doneC, stopC, err := futures.WsAllMarkPriceServeWithRate(1*time.Second, wsAllMarkPriceHandler, errHandler)
	fmt.Println(doneC, stopC)

	// use stopC to exit
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	stopC <- struct{}{}
	//}()
	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				select {
				case <-doneC:
					fmt.Println("WebSocket 已關閉，doneC停止定時寫入 Redis")
					return
				default:
					// WebSocket 仍然存活，執行寫入 Redis 的操作
					err := redisService.SetTaskHearthBeat("TaskQuote")
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done(): // STEP 2：監聽 context Done
				stopC <- struct{}{}
				fmt.Println("Task Quote cancel")
				return // kill goroutine
			}
		}
	}()
	// remove this if you do not want to be blocked here
	<-doneC
}

func backgroundTaskOrderReply(ctx context.Context, hub *wsLib.Hub, rdb *redis.Client, cancel func(), db *gorm.DB) {
	redisService := redisLib.NewRedisService(rdb)
	err := redisService.SetTaskHearthBeat("TaskOrder")
	if err != nil {
		fmt.Println(err)
	}

	orderSymbolService := orderApp.NewOrderSymbolService(db)
	leverageSymbolService := leverageApp.NewLeverageSymbolService(db)

	apiKey, secretKey := binanceLib.GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
	listenKey, err := futuresClient.NewStartUserStreamService().Do(context.Background())
	//fmt.Println("listenKey", listenKey)
	fmt.Println("Start backgroundTaskOrderReply...")
	wsUserDataHandler := func(event *futures.WsUserDataEvent) {
		//UserDataEventTypeListenKeyExpired    UserDataEventType = "listenKeyExpired"
		//UserDataEventTypeMarginCall          UserDataEventType = "MARGIN_CALL"
		//UserDataEventTypeAccountUpdate       UserDataEventType = "ACCOUNT_UPDATE"
		//UserDataEventTypeOrderTradeUpdate    UserDataEventType = "ORDER_TRADE_UPDATE"
		//UserDataEventTypeAccountConfigUpdate UserDataEventType = "ACCOUNT_CONFIG_UPDATE"
		//fmt.Println("event", event)
		userDataEventType := event.Event
		if userDataEventType == futures.UserDataEventTypeListenKeyExpired {
			err := redisService.DelTaskHearthBeat(redisLib.TASK_ORDER)
			cancel()
			log.Debug("TaskOrder listenKeyExprired")
			if err != nil {
				log.Error("TaskOrder listenKeyExprired")
			}
		} else if userDataEventType == futures.UserDataEventTypeOrderTradeUpdate {
			// 订单/交易 更新推送
			orderTradeUpdate := event.OrderTradeUpdate
			symbol := orderTradeUpdate.Symbol
			// 本次事件的具体执行类型: TRADE 交易, 订单状态: FILLED
			if orderTradeUpdate.ExecutionType == futures.OrderExecutionTypeTrade {
				if orderTradeUpdate.Status == futures.OrderStatusTypeFilled {
					fmt.Println(symbol, orderSymbolService, leverageSymbolService)
					// 已刪除 片段
					// todo: 再平衡
					//orderId := orderTradeUpdate.ID              // 訂單ID
					//tradeID := orderTradeUpdate.TradeID         // 成交ID
					//realizedPnL := orderTradeUpdate.RealizedPnL // 该交易实现盈亏
					//
					//orderSymbol, _ := orderSymbolService.GetOrderSymbolByOrderId(orderId)

				}
			}
		}
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, stopC, err := futures.WsUserDataServe(listenKey, wsUserDataHandler, errHandler)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(doneC, stopC)

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				select {
				case <-doneC:
					fmt.Println("TaskOrder WebSocket 已關閉，doneC停止定時寫入 Redis")
					return
				default:
					// WebSocket 仍然存活，執行寫入 Redis 的操作
					err := redisService.SetTaskHearthBeat("TaskOrder")
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}()
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	stopC <- struct{}{}
	//}()
	go func() {
		for {
			select {
			case <-ctx.Done(): // STEP 2：監聽 context Done
				stopC <- struct{}{}
				fmt.Println("Task Quote cancel")
				return // kill goroutine
			}
		}
	}()
	<-doneC
}

func heartBeat(hub *wsLib.Hub, rdb *redis.Client, db *gorm.DB) {
	redisService := redisLib.NewRedisService(rdb)
	var cancel func()
	var ctx context.Context
	var cancel2 func()
	var ctx2 context.Context

	ticker := time.NewTicker(1000 * time.Millisecond)

	go func() {
		//fmt.Println("Number of active goroutines 222:", runtime.NumGoroutine())
		for {
			select {
			case <-ticker.C:
				go func() {
					_, err := redisService.GetTaskHearthBeat(redisLib.TASK_QUOTE)
					//fmt.Println("heartbeat", err)
					if err != nil {
						fmt.Println("heartbeat cancel TaskQuote", cancel)
						if cancel != nil {
							cancel()
						}
						time.Sleep(1 * time.Millisecond)
						ctx, cancel = context.WithCancel(context.Background())
						go backgroundTaskQuote(ctx, hub, rdb)
					}
				}()

				go func() {
					_, err2 := redisService.GetTaskHearthBeat(redisLib.TASK_ORDER)
					//fmt.Println("heartbeat", err)
					if err2 != nil {
						fmt.Println("heartbeat cancel TaskOrder", cancel2)
						if cancel2 != nil {
							cancel2()
						}
						time.Sleep(1 * time.Millisecond)
						ctx2, cancel2 = context.WithCancel(context.Background())
						go backgroundTaskOrderReply(ctx2, hub, rdb, cancel2, db)
					}
				}()

			}
		}

	}()

}

func Init(rdb *redis.Client) {
	redisService := redisLib.RedisService{Rdb: rdb}
	err := redisService.DelTaskHearthBeat("TaskQuote")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	const config string = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s"
	pgSources := fmt.Sprintf(config,
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DATABASE"),
		os.Getenv("PG_SSLMODE"),
		os.Getenv("PG_TIME_ZONE"),
	)

	db, dbErr := gorm.Open(postgres.Open(pgSources), &gorm.Config{})

	if dbErr != nil {
		panic("failed to connect database")
	}

	//迁移 schema
	db.AutoMigrate(&activityLogApp.ActivityLog{})
	db.AutoMigrate(&leverageApp.LeverageSymbol{})
	db.AutoMigrate(&orderApp.OrderSymbol{})

	hub := wsLib.NewHub()
	go hub.Run()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	gin.SetMode(gin.DebugMode)
	Init(rdb)

	route := router.Default()
	//route.Use(cors.Default())
	route = activityLogApp.GetRoute(route, db, hub)
	route = quoteApp.GetRoute(route, db, hub, rdb)
	route = leverageApp.GetRoute(route, db, hub)
	route = mainApp.GetRoute(route, db, hub, rdb)
	route = orderApp.GetRoute(route, db, hub)

	go heartBeat(hub, rdb, db)

	fmt.Println("running local")
	log.Fatal(http.ListenAndServe(":8000", route))

}

/*
{"e":"kline","E":1682244928001,"s":"ETHUSDT","k":{"t":1682244927000,"T":1682244927999,"s":"ETHUSDT","i":"1s","f":-1,"L":-1,"o":"1870.76000000","c":"1870.76000000","h":"1870.76000000","l":"1870.76000000","v":"0.00000000","n":0,"x":true,"q":"0.00000000","V":"0.00000000","Q":"0.00000000"}}
{"e":"kline","E":1682244928001,"s":"BTCUSDT","k":{"t":1682244927000,"T":1682244927999,"s":"BTCUSDT","i":"1s","f":3091449369,"L":3091449370,"o":"27631.01000000","c":"27631.00000000","h":"27631.01000000","l":"27631.00000000","v":"0.23040000","n":2,"x":true,"q":"6366.18251000","V":"0.01100000","Q":"303.94111000"}}*/

/*

// cryptoProjectDemoPW
*/
