package leverageApp

import (
	"gorm.io/gorm"
	"time"
)

type LeverageSymbol struct {
	gorm.Model
	Symbol                        string    `json:"symbol"`          // 交易對
	MarketPrice                   string    `json:"marketPrice"`     // 市價
	InventoryVolume               string    `json:"inventoryVolume"` // 庫存量
	InventoryVolumeUpdateDateTime time.Time `json:"inventoryVolumeUpdateDateTime"`
	SettingPriceStepPercent       string    `json:"settingPriceStepPercent"` // 設定上下價格%
	SettingCalcNum                int       `json:"settingCalcNum"`          // 設定 上下筆數
	PricePrecision                int       `json:"pricePrecision"`
	QuantityPrecision             int       `json:"quantityPrecision"`
	BaseAssetPrecision            int       `json:"baseAssetPrecision"`
	QuotePrecision                int       `json:"quotePrecision"`
	FilterMinNotional             string    `json:"filterMinNotional"`
	IsSimulate                    bool      `json:"isSimulate"`
	InitMarketPrice               string    `json:"initMarketPrice"` // 一開始的市價
}
