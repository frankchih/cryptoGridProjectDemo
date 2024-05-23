package mainApp

import (
	"cryptoGridProjectDemo/redisLib"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

type Handle struct {
	DB  *gorm.DB
	Rdb *redis.Client
}
type SysCurrStatusResponse struct {
	TaskQuoteValue float64 `json:"taskQuoteValue"`
	TaskOrderValue float64 `json:"taskOrderValue"`
}

func (handle *Handle) GetSysCurrStatus(ctx *gin.Context) {
	redisService := redisLib.NewRedisService(handle.Rdb)
	sysCurrStatusResponse := SysCurrStatusResponse{}
	taskQuoteValue, _ := redisService.GetTTLTaskHearthBeat(redisLib.TASK_QUOTE)
	taskOrderValue, _ := redisService.GetTTLTaskHearthBeat(redisLib.TASK_ORDER)

	if taskQuoteValue == -2 || taskQuoteValue == -1 {
		sysCurrStatusResponse.TaskQuoteValue = float64(taskQuoteValue)
	} else {
		sysCurrStatusResponse.TaskQuoteValue = taskQuoteValue.Seconds()
	}
	if taskOrderValue == -2 || taskOrderValue == -1 {
		sysCurrStatusResponse.TaskOrderValue = float64(taskOrderValue)
	} else {
		sysCurrStatusResponse.TaskOrderValue = taskOrderValue.Seconds()
	}

	ctx.JSON(http.StatusOK, sysCurrStatusResponse)
}

type TaskRestartInput struct {
	TaskName string `json:"taskName" binding:"required,omitempty"`
}

func (handle *Handle) PostTaskRestart(ctx *gin.Context) {
	taskRestartInput := &TaskRestartInput{}
	err := ctx.ShouldBindJSON(taskRestartInput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err.Error(),
		})
		return
	}
	taskName := taskRestartInput.TaskName
	redisService := redisLib.NewRedisService(handle.Rdb)

	err2 := redisService.DelTaskHearthBeat(taskName)
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"isErr": true,
			"err":   err2.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Del success",
	})
}
