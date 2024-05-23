package orderApp

import (
	"cryptoGridProjectDemo/wsLib"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB, hub *wsLib.Hub) *gin.Engine {
	handle := &Handle{DB: db}

	routeGroup := route.Group("api/").Group("orderApp/")
	{
		routeGroup.GET("getCurrOpenOrderList/", handle.GetCurrOpenOrderList)
		routeGroup.GET("getOrderListLastThreeDay/", handle.GetOrderListLastThreeDay)
		routeGroup.POST("createFirstGridOrder/", handle.CreateFirstGridOrder)
		routeGroup.POST("cancelAllGridOrderBySymbol/", handle.CancelAllGridOrderBySymbol)
	}

	return route
}
