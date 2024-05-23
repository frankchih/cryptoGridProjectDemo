package binanceLib

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"os"
)

func GetBinanceEnv() (string, string) {
	apiKey := os.Getenv("BI_API_KEY")
	secretKey := os.Getenv("BI_SECRET_KEY")
	return apiKey, secretKey
}

func GetOrderList(symbol string) ([]*futures.Order, error) {
	// 查詢所有訂單 包含歷史訂單 要給symbol
	apiKey, secretKey := GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	orderList, err := futuresClient.NewListOrdersService().Symbol(symbol).Do(context.Background())
	return orderList, err
}
func GetOpenOrderList(symbol string) ([]*futures.Order, error) {
	// 查詢當前掛單 要給symbol
	apiKey, secretKey := GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	orderList, err := futuresClient.NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
	return orderList, err
}

func GetExchangeInfo() (*futures.ExchangeInfo, error) {
	// 查詢
	apiKey, secretKey := GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	exchangeInfo, err := futuresClient.NewExchangeInfoService().Do(context.Background())
	return exchangeInfo, err
}

/*
	GTC - Good Till Cancel 成交为止
	IOC - Immediate or Cancel 无法立即成交(吃单)的部分就撤销
	FOK - Fill or Kill 无法全部立即成交就撤销
	GTX - Good Till Crossing 无法成为挂单方就撤销
*/
func CreateBuyLimitOrder(symbol string, price string, quantity string) (*futures.CreateOrderResponse, error) {
	// 下單
	// Side 方向 Sell buy
	// type LIMIT, MARKET, STOP, TAKE_PROFIT, STOP_MARKET, TAKE_PROFIT_MARKET, TRAILING_STOP_MARKET
	// quantity 數量
	// price
	// 有效方法 timeInForce
	apiKey, secretKey := GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	resOrder, err := futuresClient.NewCreateOrderService().Symbol(symbol).Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
		Quantity(quantity).Price(price).TimeInForce(futures.TimeInForceTypeGTC).Do(context.Background())

	return resOrder, err
}

func CreateSellLimitOrder(symbol string, price string, quantity string) (*futures.CreateOrderResponse, error) {
	// 下單
	// Side 方向 Sell buy
	// type LIMIT, MARKET, STOP, TAKE_PROFIT, STOP_MARKET, TAKE_PROFIT_MARKET, TRAILING_STOP_MARKET
	// quantity 數量
	// price
	// 有效方法 timeInForce
	apiKey, secretKey := GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	resOrder, err := futuresClient.NewCreateOrderService().Symbol(symbol).Side(futures.SideTypeSell).Type(futures.OrderTypeLimit).
		Quantity(quantity).Price(price).TimeInForce(futures.TimeInForceTypeGTC).Do(context.Background())

	return resOrder, err
}

type BatchOrderInput struct {
	Price    string
	Quantity string
}

func DeleteBuyLimitBatchOrder(symbol string, orderIDList []int64) ([]*futures.CancelOrderResponse, []error) {
	// 批次取消訂單  最多支持10个订单
	apiKey, secretKey := GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	batchSize := 10
	cancelOrderResponseList := make([]*futures.CancelOrderResponse, 0)
	errList := make([]error, 0)
	for i := 0; i < len(orderIDList); i += batchSize {
		end := i + batchSize
		if end > len(orderIDList) {
			end = len(orderIDList)
		}
		batchOrderIDList := orderIDList[i:end]

		cancelOrderResponse, err := futuresClient.NewCancelMultipleOrdersService().Symbol(symbol).OrderIDList(batchOrderIDList).Do(context.Background())
		cancelOrderResponseList = append(cancelOrderResponseList, cancelOrderResponse...)
		errList = append(errList, err)
	}

	return cancelOrderResponseList, errList
}

func Test() {
	//resOrder, err := CreateLimitOrder("AVAXUSDT", futures.SideTypeBuy, "1", "10")
	//fmt.Println(resOrder, err)

	openOrderList, err := GetOpenOrderList("AVAXUSDT")
	fmt.Println(openOrderList, err)

	orderList, err := GetOrderList("AVAXUSDT")
	fmt.Println(orderList, err)
}
