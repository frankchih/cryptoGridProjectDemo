package activityLogApp

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Handle struct {
	DB *gorm.DB
}

func (handle *Handle) Main(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Main"})
}
func (handle *Handle) Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Hello World"})
}

//type BookmarkSaveFromChromeData struct {
//	Bookmarks []*Bookmark `json:"bookmarks" binding:"required,omitempty"`
//}

func (handle *Handle) CreateTestActivityLog(ctx *gin.Context) {

	err := CreateActivityLogTest(handle.DB, "test msg")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (handle *Handle) GetActivityLogs(ctx *gin.Context) {
	db := handle.DB
	activityLogService := NewActivityLogService(db)
	activityLogs, err := activityLogService.GetActivityLogs()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "activityLogs": activityLogs})
}
