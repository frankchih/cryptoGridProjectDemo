package quoteApp

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

type Handle struct {
	DB  *gorm.DB
	Rdb *redis.Client
}

func (handle *Handle) Main(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Hi "})
}
