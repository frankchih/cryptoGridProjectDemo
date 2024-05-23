package quoteApp

import (
	"cryptoGridProjectDemo/wsLib"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetRoute(route *gin.Engine, db *gorm.DB, hub *wsLib.Hub, rdb *redis.Client) *gin.Engine {
	handle := &Handle{DB: db, Rdb: rdb}
	fmt.Println(handle)
	route.GET("ws/quote/", func(c *gin.Context) {
		wsLib.ServeWs(hub, c.Writer, c.Request, "ws/quote/")
	})

	return route
}
