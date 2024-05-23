package activityLogApp

import (
	"cryptoGridProjectDemo/internal/pkg/util"
	"gorm.io/gorm"
)

type ActivityLogDao struct {
	db *gorm.DB
}

func NewActivityLogDao(db *gorm.DB) *ActivityLogDao {
	return &ActivityLogDao{db: db}
}

func (activityLogDao *ActivityLogDao) CreateInfo(message string) error {
	activityLog := &ActivityLog{Uuid: util.GenerateUUID(), Category: "INFO", Message: message}
	activityLogService := ActivityLogService{db: activityLogDao.db}
	return activityLogService.CreateActivityLog(activityLog)
}
