package orderApp

import (
	"cryptoGridProjectDemo/binanceLib"
	"cryptoGridProjectDemo/leverageApp"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type Handle struct {
	DB *gorm.DB
}

func (handle *Handle) GetOrderSymbolList(ctx *gin.Context) {
	orderSymbolService := NewOrderSymbolService(handle.DB)
	orderSymbols, err := orderSymbolService.GetOrderSymbols()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "orderSymbols": orderSymbols})
}

func order2orderSymbol(order *futures.Order) OrderSymbol {
	/*
		畢安的回復order 轉成 OrderSymbol
	*/
	return OrderSymbol{
		OrderId:       order.OrderID,
		Symbol:        order.Symbol,
		ClientOrderId: order.ClientOrderID,
		CumQuote:      order.CumQuote,
		ExecutedQty:   order.ExecutedQuantity,
		AvgPrice:      order.AvgPrice,
		Price:         order.Price,
		ReduceOnly:    order.ReduceOnly,
		UpdateTime:    order.UpdateTime,
		Side:          order.Side,
		Type:          order.Type,
		TimeInForce:   order.TimeInForce,
		Status:        order.Status,
		PositionSide:  order.PositionSide,
	}
}
func createOrderResponse2orderSymbol(createOrderResponse *futures.CreateOrderResponse) OrderSymbol {
	/*
		畢安的回復createOrderResponse 轉成 OrderSymbol
	*/
	return OrderSymbol{
		OrderId:       createOrderResponse.OrderID,
		Symbol:        createOrderResponse.Symbol,
		ClientOrderId: createOrderResponse.ClientOrderID,
		CumQuote:      createOrderResponse.CumQuote,
		ExecutedQty:   createOrderResponse.ExecutedQuantity,
		AvgPrice:      createOrderResponse.AvgPrice,
		Price:         createOrderResponse.Price,
		ReduceOnly:    createOrderResponse.ReduceOnly,
		UpdateTime:    createOrderResponse.UpdateTime,
		Side:          createOrderResponse.Side,
		Type:          createOrderResponse.Type,
		TimeInForce:   createOrderResponse.TimeInForce,
		Status:        createOrderResponse.Status,
		PositionSide:  createOrderResponse.PositionSide,
	}
}

//OrigQuantity      string           `json:"origQty"`
//StopPrice         string           `json:"stopPrice"`
//WorkingType       WorkingType      `json:"workingType"`
//ActivatePrice     string           `json:"activatePrice"`
//PriceRate         string           `json:"priceRate"`
//ClosePosition     bool             `json:"closePosition"`
//PriceProtect      bool             `json:"priceProtect"`
//RateLimitOrder10s string           `json:"rateLimitOrder10s,omitempty"`
//RateLimitOrder1m  string           `json:"rateLimitOrder1m,omitempty"`

func cancelOrderResponse2orderSymbol(cancelOrderResponse *futures.CancelOrderResponse) OrderSymbol {
	/*
		畢安的 cancelOrderResponse 轉成 OrderSymbol
	*/
	//cancelOrderResponse.CumQuantity 沒解釋
	//cancelOrderResponse.OrigQuantity
	//cancelOrderResponse StopPrice
	//WorkingType
	//ActivatePrice
	//PriceRate
	//OrigType
	//PriceProtect
	return OrderSymbol{
		OrderId:       cancelOrderResponse.OrderID,
		Symbol:        cancelOrderResponse.Symbol,
		ClientOrderId: cancelOrderResponse.ClientOrderID,
		CumQuote:      cancelOrderResponse.CumQuote,
		ExecutedQty:   cancelOrderResponse.ExecutedQuantity,
		//AvgPrice:      order.AvgPrice,
		Price:        cancelOrderResponse.Price,
		ReduceOnly:   cancelOrderResponse.ReduceOnly,
		UpdateTime:   cancelOrderResponse.UpdateTime,
		Side:         cancelOrderResponse.Side,
		Type:         cancelOrderResponse.Type,
		TimeInForce:  cancelOrderResponse.TimeInForce,
		Status:       cancelOrderResponse.Status,
		PositionSide: cancelOrderResponse.PositionSide,
	}
}

func (handle *Handle) GetCurrOpenOrderList(ctx *gin.Context) {
	/*
		抓當前開單 並更新資料庫，要改成 bulk update 和bulk create
	*/
	symbol := ctx.Query("symbol")

	openOrderList, err := binanceLib.GetOpenOrderList(symbol)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderIdList []int64
	orderSymbolService := NewOrderSymbolService(handle.DB)
	for _, openOrder := range openOrderList {
		orderId := openOrder.OrderID
		orderSymbol := order2orderSymbol(openOrder)
		handle.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "order_id"}}, // key column
			DoUpdates: clause.AssignmentColumns([]string{"symbol",
				"client_order_id", "cum_quote", "executed_qty", "avg_price", "price", "reduce_only",
				"update_time", "side", "type", "time_in_force", "status", "position_side"}), // column needed to be updated
		}).Create(&orderSymbol)
		orderIdList = append(orderIdList, orderId)
	}
	orderSymbols, err := orderSymbolService.GetOrderSymbolByOrderIdList(orderIdList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "orderSymbols": orderSymbols})
}

func (handle *Handle) GetOrderListLastThreeDay(ctx *gin.Context) {
	/*
		抓當前訂單 並更新資料庫，要改成 bulk update 和bulk create
	*/
	symbol := ctx.Query("symbol")

	openOrderList, err := binanceLib.GetOrderList(symbol)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderIdList []int64
	orderSymbolService := NewOrderSymbolService(handle.DB)
	for _, openOrder := range openOrderList {
		orderId := openOrder.OrderID
		orderSymbol := order2orderSymbol(openOrder)
		handle.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "order_id"}}, // key column
			DoUpdates: clause.AssignmentColumns([]string{"symbol",
				"client_order_id", "cum_quote", "executed_qty", "avg_price", "price", "reduce_only",
				"update_time", "side", "type", "time_in_force", "status", "position_side"}), // column needed to be updated
		}).Create(&orderSymbol)
		orderIdList = append(orderIdList, orderId)
	}
	orderSymbols, err := orderSymbolService.GetOrderSymbolByOrderIdList(orderIdList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "orderSymbols": orderSymbols})
}

type CreateFirstGridOrderInput struct {
	LeverageSymbolId        int    `json:"leverageSymbolId" binding:"required"`
	SettingPriceStepPercent string `json:"settingPriceStepPercent" binding:"required"` // 設定上下價格%
	SettingCalcNum          int    `json:"settingCalcNum" binding:"required"`          // 設定 上下筆數
	InventoryVolume         string `json:"inventoryVolume" binding:"required"`         // 庫存量
	MarketPrice             string `json:"marketPrice" binding:"required"`             // 市價
	IsSimulate              bool   `json:"isSimulate"`                                 // 是否模擬單
}

func (handle *Handle) CreateFirstGridOrder(ctx *gin.Context) {
	/*
		第一次下單 ，會產生上單下單
	*/
	createFirstGridOrderInput := &CreateFirstGridOrderInput{}
	err := ctx.ShouldBindJSON(createFirstGridOrderInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}
	leverageSymbolId := createFirstGridOrderInput.LeverageSymbolId
	settingPriceStepPercent := createFirstGridOrderInput.SettingPriceStepPercent
	settingCalcNum := createFirstGridOrderInput.SettingCalcNum
	inventoryVolume := createFirstGridOrderInput.InventoryVolume
	marketPrice := createFirstGridOrderInput.MarketPrice
	isSimulate := true // 先固定設true
	// todo: 要不要先取消所有掛單

	leverageSymbolService := leverageApp.NewLeverageSymbolService(handle.DB)
	leverageSymbol, err := leverageSymbolService.GetLeverageSymbol(leverageSymbolId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}
	handle.DB.Model(&leverageSymbol).Updates(map[string]interface{}{
		"setting_price_step_percent": settingPriceStepPercent,
		"setting_calc_num":           settingCalcNum,
		"inventory_volume":           inventoryVolume,
		"market_price":               marketPrice,
		"is_simulate":                isSimulate,
		"init_market_price":          marketPrice,
	})
	symbol := leverageSymbol.Symbol

	pricePrecision := leverageSymbol.PricePrecision

	// 網格下單
	calcLimitOrderExceptDownList := CalcLimitOrderExceptDown(settingCalcNum, settingPriceStepPercent, marketPrice, marketPrice, pricePrecision)
	batchOrderInputDownList := make([]*binanceLib.BatchOrderInput, 0)
	for _, calcLimitOrderExceptDown := range calcLimitOrderExceptDownList {
		temp := &binanceLib.BatchOrderInput{
			Price:    calcLimitOrderExceptDown.ExceptPrice,
			Quantity: calcLimitOrderExceptDown.DiffExpectAssetVolume,
		}
		batchOrderInputDownList = append(batchOrderInputDownList, temp)
	}

	//orderSymbolService := NewOrderSymbolService(handle.DB)
	orderSymbolDownList := make([]*OrderSymbol, 0)
	//orderIdDownList := make([]int64, 0)

	if !isSimulate {
		// 刪除片段
		//resOrderDownList, errList := binanceLib.CreateBuyLimitBatchOrder(symbol, batchOrderInputDownList)

	} else {
		for _, batchOrderInputDown := range batchOrderInputDownList {
			orderSymbol := OrderSymbol{
				Symbol:      symbol,
				Price:       batchOrderInputDown.Price,
				ExecutedQty: batchOrderInputDown.Quantity,
			}
			orderSymbolDownList = append(orderSymbolDownList, &orderSymbol)
		}
		handle.DB.Clauses(clause.Returning{}).Create(&orderSymbolDownList)
	}
	//fmt.Println("up start")
	// 網格上單
	calcLimitOrderExceptUpList := CalcLimitOrderExceptUp(settingCalcNum, settingPriceStepPercent, marketPrice, marketPrice, pricePrecision)
	batchOrderInputUpList := make([]*binanceLib.BatchOrderInput, 0)
	for _, calcLimitOrderExceptUp := range calcLimitOrderExceptUpList {
		temp := &binanceLib.BatchOrderInput{
			Price:    calcLimitOrderExceptUp.ExceptPrice,
			Quantity: calcLimitOrderExceptUp.DiffExpectAssetVolume,
		}
		batchOrderInputUpList = append(batchOrderInputUpList, temp)
	}

	orderSymbolUpList := make([]*OrderSymbol, 0)
	//orderIdUpList := make([]int64, 0)

	if !isSimulate {
		// 刪除片段
		//resOrderUpList, errList := binanceLib.CreateSellLimitBatchOrder(symbol, batchOrderInputUpList)
	} else {
		for _, batchOrderInputUp := range batchOrderInputUpList {
			orderSymbol := OrderSymbol{
				Symbol:      symbol,
				Price:       batchOrderInputUp.Price,
				ExecutedQty: batchOrderInputUp.Quantity,
			}
			orderSymbolUpList = append(orderSymbolUpList, &orderSymbol)
		}
		handle.DB.Clauses(clause.Returning{}).Create(&orderSymbolUpList)
	}

	//orderSymbolDowns, err := orderSymbolService.GetOrderSymbolByOrderIdList(orderIdDownList)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//orderSymbolUps, err := orderSymbolService.GetOrderSymbolByOrderIdList(orderIdUpList)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":                 "success",
		"orderSymbolDownList": orderSymbolDownList,
		"orderSymbolUpList":   orderSymbolUpList,
	})

}

type CancelAllGridOrderBySymbolInput struct {
	LeverageSymbolId int `json:"leverageSymbolId" binding:"required"`
}

func (handle *Handle) CancelAllGridOrderBySymbol(ctx *gin.Context) {
	cancelAllGridOrderBySymbolInput := &CancelAllGridOrderBySymbolInput{}
	err := ctx.ShouldBindJSON(cancelAllGridOrderBySymbolInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}
	leverageSymbolId := cancelAllGridOrderBySymbolInput.LeverageSymbolId

	leverageSymbolService := leverageApp.NewLeverageSymbolService(handle.DB)
	leverageSymbol, err := leverageSymbolService.GetLeverageSymbol(leverageSymbolId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}
	symbol := leverageSymbol.Symbol

	// 先抓開單
	openOrderList, err := binanceLib.GetOpenOrderList(symbol)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderIdList := make([]int64, 0)
	orderIdMapOrderSymbol := map[int64]*OrderSymbol{}
	orderSymbolUpdateList := make([]*OrderSymbol, 0)

	for _, openOrder := range openOrderList {
		orderSymbol := order2orderSymbol(openOrder)
		orderIdMapOrderSymbol[openOrder.OrderID] = &orderSymbol
		orderIdList = append(orderIdList, openOrder.OrderID)

		orderSymbolUpdateList = append(orderSymbolUpdateList, &orderSymbol)
	}

	orderSymbolService := NewOrderSymbolService(handle.DB)
	//先查資料庫 有甚麼資料
	orderSymbols, err := orderSymbolService.GetOrderSymbolByOrderIdList(orderIdList)
	//刪除 map 的資料
	for _, orderSymbol := range orderSymbols {
		delete(orderIdMapOrderSymbol, orderSymbol.OrderId)
	}
	//就會剩下需要建立的 orderSymbol 資料
	orderSymbolCreateList := make([]*OrderSymbol, 0)
	for _, orderSymbol := range orderIdMapOrderSymbol {
		orderSymbolCreateList = append(orderSymbolCreateList, orderSymbol)
	}
	handle.DB.Create(orderSymbolCreateList)

	//直接更新全部
	handle.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "order_id"}}, // key column
		DoUpdates: clause.AssignmentColumns([]string{"symbol",
			"client_order_id", "cum_quote", "executed_qty", "avg_price", "price", "reduce_only",
			"update_time", "side", "type", "time_in_force", "status", "position_side"}), // column needed to be updated
	}).Updates(&orderSymbolUpdateList)

	// 刪除片段
	//cancelOrderResponseList, errList := binanceLib.DeleteBuyLimitBatchOrder(symbol, orderIdList)
	//
	//for _, err := range errList {
	//	if err != nil {
	//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//		return
	//	}
	//}
	//cancelOrderSymbolUpdateList := make([]*OrderSymbol, 0)

	//for _, cancelOrderResponse := range cancelOrderResponseList {
	//	orderSymbol := cancelOrderResponse2orderSymbol(cancelOrderResponse)
	//	cancelOrderSymbolUpdateList = append(cancelOrderSymbolUpdateList, &orderSymbol)
	//}

	////直接更新全部
	//handle.DB.Clauses(clause.OnConflict{
	//	Columns: []clause.Column{{Name: "order_id"}}, // key column
	//	DoUpdates: clause.AssignmentColumns([]string{"symbol",
	//		"client_order_id", "cum_quote", "executed_qty", "avg_price", "price", "reduce_only",
	//		"update_time", "side", "type", "time_in_force", "status", "position_side"}), // column needed to be updated
	//}).Updates(&cancelOrderSymbolUpdateList)

	orderSymbolResultList, err := orderSymbolService.GetOrderSymbolByOrderIdList(orderIdList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "orderSymbols": orderSymbolResultList})
}

func (handle *Handle) Test(ctx *gin.Context) {
	//resOrder, err := binanceLib.CreateBuyLimitOrder("AVAXUSDT", "10", "1")
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//fmt.Println(resOrder)
	//resOrder2, err := binanceLib.CreateSellLimitOrder("AVAXUSDT", "40", "1")
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//fmt.Println(resOrder2)
	//ctx.JSON(http.StatusOK, gin.H{"msg": "success", "resOrder": resOrder})
}

type CalcLimitOrderExcept struct {
	SeqNum                int    // 序列號
	ExceptPrice           string // 預期價格 下單價格
	DiffExpectAssetVolume string // 預期增加/減少量 下單量
}

func CalcLimitOrderExceptDown(calcNum int, priceStep, marketPrice, initPrice string, pricePrecision int) []*CalcLimitOrderExcept {
	/*
		計算往 "下" 單 ，刪除片段，僅留下 價錢 計算，數量直接壓0
	*/
	priceStepDecimal, _ := decimal.NewFromString(priceStep)
	marketPriceDecimal, _ := decimal.NewFromString(marketPrice)
	initPriceDecimal, _ := decimal.NewFromString(initPrice)

	dynamicMarketPriceDecimal := marketPriceDecimal
	calcLimitOrderExceptList := make([]*CalcLimitOrderExcept, calcNum)

	fixDiffPrice := initPriceDecimal.Mul(priceStepDecimal) // 固定差額

	for i := 0; i < calcNum; i++ {
		exceptPrice := dynamicMarketPriceDecimal.Sub(fixDiffPrice) // 預期價格 = 市價 - 固定差額
		exceptPrice = exceptPrice.RoundDown(int32(pricePrecision)) // 無條件捨去

		calcLimitOrderExceptList[i] = &CalcLimitOrderExcept{
			SeqNum:                i,
			ExceptPrice:           exceptPrice.String(),
			DiffExpectAssetVolume: "0",
		}
		dynamicMarketPriceDecimal = exceptPrice
	}
	return calcLimitOrderExceptList
}

func CalcLimitOrderExceptUp(calcNum int, priceStep, marketPrice, initPrice string, pricePrecision int) []*CalcLimitOrderExcept {
	/*
		計算往 "上" 單，刪除片段，僅留下 價錢 計算，數量直接壓0
	*/
	priceStepDecimal, _ := decimal.NewFromString(priceStep)
	marketPriceDecimal, _ := decimal.NewFromString(marketPrice)
	initPriceDecimal, _ := decimal.NewFromString(initPrice)

	dynamicMarketPriceDecimal := marketPriceDecimal
	calcLimitOrderExceptList := make([]*CalcLimitOrderExcept, calcNum)

	fixDiffPrice := initPriceDecimal.Mul(priceStepDecimal) // 固定差額

	for i := 0; i < calcNum; i++ {
		exceptPrice := dynamicMarketPriceDecimal.Add(fixDiffPrice) // 預期價格 = 市價 + 固定差額
		exceptPrice = exceptPrice.RoundDown(int32(pricePrecision)) // 無條件捨去

		calcLimitOrderExceptList[i] = &CalcLimitOrderExcept{
			SeqNum:                i,
			ExceptPrice:           exceptPrice.String(),
			DiffExpectAssetVolume: "0",
		}
		dynamicMarketPriceDecimal = exceptPrice
	}

	return calcLimitOrderExceptList
}

/*
{"symbol":"AVAXUSDT",
"orderId":14976512647,
"clientOrderId":"XtxGGQhhyMUQfhEHf5e504",
"price":"10",
"origQty":"1",
"executedQty":"0",
"cumQuote":"0",
"reduceOnly":false,
"status":"NEW",
"stopPrice":"0",
"timeInForce":"GTC",
"type":"LIMIT",
"side":"BUY",
"updateTime":1685782976536,
"workingType":"CONTRACT_PRICE",
"activatePrice":"","priceRate":"","avgPrice":"0.0000","positionSide":"BOTH","closePosition":false,"priceProtect":false,"rateLimitOrder10s":"0","rateLimitOrder1m":"1"}}
//https://stackoverflow.com/questions/39333102/how-to-create-or-update-a-record-with-gorm
// https://binance-docs.github.io/apidocs/futures/cn/#user_data
*/
