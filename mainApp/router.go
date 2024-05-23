package mainApp

import (
	"cryptoGridProjectDemo/wsLib"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB, hub *wsLib.Hub, rdb *redis.Client) *gin.Engine {
	handle := &Handle{DB: db, Rdb: rdb}

	routeGroup := route.Group("api/").Group("mainApp/")
	{
		routeGroup.GET("getSysCurrStatus/", handle.GetSysCurrStatus)
		routeGroup.POST("postTaskRestart/", handle.PostTaskRestart)
	}

	return route
}
