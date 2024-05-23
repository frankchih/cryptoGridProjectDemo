package orderApp

import (
	"github.com/adshao/go-binance/v2/futures"
	"gorm.io/gorm"
)

type OrderSymbol struct {
	gorm.Model
	Symbol               string                     `json:"symbol"`               // 交易对
	OrderId              int64                      `json:"orderId"`              // 系统订单号
	ClientOrderId        string                     `json:"clientOrderId"`        // 用户自定义的订单号
	CumQuote             string                     `json:"cumQuote"`             // 成交金额
	ExecutedQty          string                     `json:"executedQty"`          // 成交量
	AvgPrice             string                     `json:"avgPrice"`             // 平均成交价
	Price                string                     `json:"price"`                // 委托价格
	ReduceOnly           bool                       `json:"reduceOnly"`           // 仅减仓
	UpdateTime           int64                      `json:"updateTime"`           // 更新时间 --以上訂單
	Side                 futures.SideType           `json:"side"`                 // 订单方向 SideType
	Type                 futures.OrderType          `json:"type"`                 // 订单类型 OrderType
	TimeInForce          futures.TimeInForceType    `json:"timeInForce"`          // 有效方式 TimeInForceType GTC IOC FOK GTX
	OriginalQty          string                     `json:"originalQty"`          // 订单原始数量
	OriginalPrice        string                     `json:"originalPrice"`        // 订单原始价格
	AveragePrice         string                     `json:"averagePrice"`         // 订单平均价格
	StopPrice            string                     `json:"stopPrice"`            // 条件订单触发价格，对追踪止损单无效
	ExecutionType        futures.OrderExecutionType `json:"executionType"`        // 本次事件的具体执行类型 OrderExecutionType
	Status               futures.OrderStatusType    `json:"status"`               // 订单的当前状态 OrderStatusType
	LastFilledQty        string                     `json:"lastFilledQty"`        // 订单末次成交量
	AccumulatedFilledQty string                     `json:"accumulatedFilledQty"` // 订单累计已成交量
	LastFilledPrice      string                     `json:"lastFilledPrice"`      // 订单末次成交价格
	CommissionAsset      string                     `json:"commissionAsset"`      // 手续费资产类型
	Commission           string                     `json:"commission"`           // 手续费数量
	TradeTime            int64                      `json:"tradeTime"`            // 成交时间
	TradeID              int64                      `json:"tradeID"`              // 成交ID
	BidsNotional         string                     `json:"bidsNotional"`         // 买单净值
	AsksNotional         string                     `json:"asksNotional"`         // 卖单净值
	IsMaker              bool                       `json:"isMaker"`              // 该成交是作为挂单成交吗？
	IsReduceOnly         bool                       `json:"isReduceOnly"`         // 是否是只减仓单
	WorkingType          futures.WorkingType        `json:"workingType"`          // 触发价类型 WorkingType
	OriginalType         futures.OrderType          `json:"originalType"`         // 原始订单类型 OrderType
	PositionSide         futures.PositionSideType   `json:"positionSide"`         // 持仓方向 PositionSideType "BOTH" "LONG" "SHORT"
	IsClosingPosition    bool                       `json:"isClosingPosition"`    // 是否为触发平仓单; 仅在条件订单情况下会推送此字段
	ActivationPrice      string                     `json:"activationPrice"`      // 追踪止损激活价格, 仅在追踪止损单时会推送此字段
	CallbackRate         string                     `json:"callbackRate"`         // 追踪止损回调比例, 仅在追踪止损单时会推送此字段
	RealizedPnL          string                     `json:"realizedPnL"`          // 该交易实现盈亏
}

/*
SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"

	PositionSideTypeBoth  PositionSideType = "BOTH"
	PositionSideTypeLong  PositionSideType = "LONG"
	PositionSideTypeShort PositionSideType = "SHORT"

	OrderTypeLimit              OrderType = "LIMIT"
	OrderTypeMarket             OrderType = "MARKET"
	OrderTypeStop               OrderType = "STOP"
	OrderTypeStopMarket         OrderType = "STOP_MARKET"
	OrderTypeTakeProfit         OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket   OrderType = "TAKE_PROFIT_MARKET"
	OrderTypeTrailingStopMarket OrderType = "TRAILING_STOP_MARKET"

	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceTypeIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceTypeGTX TimeInForceType = "GTX" // Good Till Crossing (Post Only)

   "S":"SELL",                     // 订单方向
   "o":"TRAILING_STOP_MARKET", // 订单类型
   "f":"GTC",                      // 有效方式
   "q":"0.001",                    // 订单原始数量
   "p":"0",                        // 订单原始价格
   "ap":"0",                       // 订单平均价格
   "sp":"7103.04",                 // 条件订单触发价格，对追踪止损单无效
   "x":"NEW",                      // 本次事件的具体执行类型
   "X":"NEW",                      // 订单的当前状态
   "i":8886774,                    // 订单ID
   "l":"0",                        // 订单末次成交量
   "z":"0",                        // 订单累计已成交量
   "L":"0",                        // 订单末次成交价格
   "N": "USDT",                    // 手续费资产类型
   "n": "0",                       // 手续费数量
   "T":1568879465650,              // 成交时间
   "t":0,                          // 成交ID
   "b":"0",                        // 买单净值
   "a":"9.91",                     // 卖单净值
   "m": false,                     // 该成交是作为挂单成交吗？
   "R":false   ,                   // 是否是只减仓单
   "wt": "CONTRACT_PRICE",         // 触发价类型
   "ot": "TRAILING_STOP_MARKET",   // 原始订单类型
   "ps":"LONG"                     // 持仓方向
   "cp":false,                     // 是否为触发平仓单; 仅在条件订单情况下会推送此字段
   "AP":"7476.89",                 // 追踪止损激活价格, 仅在追踪止损单时会推送此字段
   "cr":"5.0",                     // 追踪止损回调比例, 仅在追踪止损单时会推送此字段
   "pP": false,              // 忽略
   "si": 0,                  // 忽略
   "ss": 0,                  // 忽略
   "rp":"0"                       // 该交易实现盈亏

	order: 欄位
	ClientOrderId string `json:"clientOrderId"`
	CumQuote      string `json:"cumQuote"`
	ExecutedQty   string `json:"executedQty"`
	OrderId       int    `json:"orderId"`
	AvgPrice      string `json:"avgPrice"`
	OrigQty       string `json:"origQty"`
	Price         string `json:"price"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	ClosePosition bool   `json:"closePosition"`
	Symbol        string `json:"symbol"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	OrigType      string `json:"origType"`
	ActivatePrice string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	UpdateTime    int64  `json:"updateTime"`
	WorkingType   string `json:"workingType"`
	PriceProtect  bool   `json:"priceProtect"`
*/
