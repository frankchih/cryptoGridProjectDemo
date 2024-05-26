package main

import (
	"context"
	"cryptoGridProjectDemo/binanceLib"
	"cryptoGridProjectDemo/internal/pkg/log"
	"cryptoGridProjectDemo/leverageApp"
	"cryptoGridProjectDemo/orderApp"
	"cryptoGridProjectDemo/redisLib"
	"cryptoGridProjectDemo/wsLib"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

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
