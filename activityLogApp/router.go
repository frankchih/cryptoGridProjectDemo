package activityLogApp

import (
	"cryptoGridProjectDemo/wsLib"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB, hub *wsLib.Hub) *gin.Engine {
	handle := &Handle{DB: db}
	route.GET("", handle.Main)

	routeGroup := route.Group("api/activityLog/")
	{
		routeGroup.GET("hello", handle.Hello)
		routeGroup.GET("getList/", handle.GetActivityLogs)
		routeGroup.GET("createTestActivityLog", handle.CreateTestActivityLog)
	}

	return route
}
