package leverageApp

import (
	"context"
	"cryptoGridProjectDemo/binanceLib"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Handle struct {
	DB *gorm.DB
}

func (handle *Handle) GetLeverageSymbolList(ctx *gin.Context) {
	leverageSymbolService := NewLeverageSymbolService(handle.DB)
	leverageSymbols, err := leverageSymbolService.GetLeverageSymbols()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "leverageSymbols": leverageSymbols})
}

type LeverageSymbolCreateInput struct {
	Symbol string `json:"symbol" binding:"required"`
	//StopPrice    string `json:"stopPrice" binding:"required"`
	//ProfitPrice  string `json:"profitPrice" binding:"required"`
	//Direction    string `json:"direction" binding:"required"` // LONG SHORT
	//MaxLossMoney string `json:"maxLossMoney" binding:"required"`
}

func (handle *Handle) CreateLeverageSymbol(ctx *gin.Context) {
	leverageSymbolInput := &LeverageSymbolCreateInput{}
	err := ctx.ShouldBindJSON(leverageSymbolInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}
	leverageSymbolService := NewLeverageSymbolService(handle.DB)
	leverageSymbol := &LeverageSymbol{Symbol: leverageSymbolInput.Symbol}
	err2 := leverageSymbolService.CreateLeverageSymbol(leverageSymbol)
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"isErr": true, "error": err2.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "leverageSymbol": leverageSymbol})
}

type LeverageSymbolUpdateInput struct {
	InventoryVolume string `json:"inventoryVolume"` // 庫存量
}

func (handle *Handle) UpdateLeverageSymbol(ctx *gin.Context) {
	leverageSymbolUpdateInput := &LeverageSymbolUpdateInput{}
	err := ctx.ShouldBindJSON(leverageSymbolUpdateInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}
	leverageSymbolId := ctx.Param("leverageSymbolId")
	leverageSymbolIdInt, _ := strconv.Atoi(leverageSymbolId)
	leverageSymbolService := NewLeverageSymbolService(handle.DB)
	leverageSymbol, err := leverageSymbolService.GetLeverageSymbol(leverageSymbolIdInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "leverageSymbol": leverageSymbol})
}

func (handle *Handle) DeleteLeverageSymbol(ctx *gin.Context) {
	leverageSymbolId := ctx.Param("leverageSymbolId")
	leverageSymbolIdInt, _ := strconv.Atoi(leverageSymbolId)
	leverageSymbolService := NewLeverageSymbolService(handle.DB)
	err := leverageSymbolService.DeleteLeverageSymbol(leverageSymbolIdInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"isErr": true, "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (handle *Handle) SendLeverageSymbol(ctx *gin.Context) {
	leverageSymbolId := ctx.Param("leverageSymbolId")
	leverageSymbolIdInt, _ := strconv.Atoi(leverageSymbolId)
	fmt.Println(leverageSymbolIdInt)

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

type ExchangeInfoDetail struct {
	Symbol             string `json:"symbol"`
	PricePrecision     int    `json:"pricePrecision"`
	QuantityPrecision  int    `json:"quantityPrecision"`
	BaseAssetPrecision int    `json:"baseAssetPrecision"`
	QuotePrecision     int    `json:"quotePrecision"`
	FilterMinNotional  string `json:"filterMinNotional"`
}

func (handle *Handle) GetCurrAsset(ctx *gin.Context) {
	/*
		抓當前資產，會更新 InventoryVolume
	*/
	apiKey, secretKey := binanceLib.GetBinanceEnv()
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	account, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"isErr": true, "error": err.Error()})
		return
	}

	//assets := account.Assets
	positions := account.Positions
	//for _, asset := range assets {
	//	fmt.Println(asset)
	//}
	hashMap := map[string]*futures.AccountPosition{}
	for _, position := range positions {
		symbol := position.Symbol
		hashMap[symbol] = position
	}

	hashMapExchangeInfoDetail := map[string]*ExchangeInfoDetail{}

	exchangeInfo, err := binanceLib.GetExchangeInfo()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"isErr": true, "error": err.Error()})
		return
	}

	//XLMUSDT 5 0 8 8
	for _, symbolObj := range exchangeInfo.Symbols {
		symbol := symbolObj.Symbol
		//FilterMinNotional

		filterMinNotional := ""
		for _, filters := range symbolObj.Filters {
			if filters["filterType"] == "MIN_NOTIONAL" {
				filterMinNotional = filters["notional"].(string)
				break
			}
		}

		hashMapExchangeInfoDetail[symbol] = &ExchangeInfoDetail{
			Symbol:             symbolObj.Symbol,
			PricePrecision:     symbolObj.PricePrecision,
			QuantityPrecision:  symbolObj.QuantityPrecision,
			BaseAssetPrecision: symbolObj.BaseAssetPrecision,
			QuotePrecision:     symbolObj.QuotePrecision,
			FilterMinNotional:  filterMinNotional,
		}
		//fmt.Println(symbolObj.Symbol, symbolObj.PricePrecision, symbolObj.QuantityPrecision, symbolObj.BaseAssetPrecision, symbolObj.QuotePrecision)
		//  "pricePrecision": 5,  // 价格小数点位数(仅作为系统精度使用，注意同tickSize 区分）
		//	"quantityPrecision": 0,  // 数量小数点位数(仅作为系统精度使用，注意同stepSize 区分）
		//	"baseAssetPrecision": 8,  // 标的资产精度
		//	"quotePrecision": 8,  // 报价资产精度
	}

	leverageSymbolService := NewLeverageSymbolService(handle.DB)
	leverageSymbols, err := leverageSymbolService.GetLeverageSymbols()
	if err != nil {

	}
	for _, leverageSymbol := range *leverageSymbols {
		symbol := leverageSymbol.Symbol
		if _, ok := hashMap[symbol]; ok {
			position := hashMap[symbol]
			positionAmt := position.PositionAmt

			handle.DB.Model(&leverageSymbol).Updates(map[string]interface{}{"inventory_volume": positionAmt, "inventory_volume_update_date_time": time.Now()})
		}

		if _, ok := hashMapExchangeInfoDetail[symbol]; ok {
			exchangeInfoDetail := hashMapExchangeInfoDetail[symbol]

			updateHashMap := map[string]interface{}{
				"price_precision":      exchangeInfoDetail.PricePrecision,
				"quantity_precision":   exchangeInfoDetail.QuantityPrecision,
				"base_asset_precision": exchangeInfoDetail.BaseAssetPrecision,
				"quote_precision":      exchangeInfoDetail.QuotePrecision,
				"filter_min_notional":  exchangeInfoDetail.FilterMinNotional,
			}
			handle.DB.Model(&leverageSymbol).Updates(updateHashMap)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}
