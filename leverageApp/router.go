package leverageApp

import (
	"cryptoGridProjectDemo/wsLib"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB, hub *wsLib.Hub) *gin.Engine {
	handle := &Handle{DB: db}

	routeGroup := route.Group("api/").Group("leverageApp/")
	{
		routeGroup.GET("getLeverageSymbolList/", handle.GetLeverageSymbolList)
		routeGroup.POST("createLeverageSymbol/", handle.CreateLeverageSymbol)
		routeGroup.PATCH("updateLeverageSymbol/:leverageSymbolId/", handle.UpdateLeverageSymbol)
		routeGroup.DELETE("deleteLeverageSymbol/:leverageSymbolId/", handle.DeleteLeverageSymbol)
		routeGroup.POST("sendLeverageSymbol/:leverageSymbolId/", handle.SendLeverageSymbol)
		routeGroup.POST("getCurrAsset/", handle.GetCurrAsset)
	}

	return route
}
