package activityLogApp

import (
	"cryptoGridProjectDemo/internal/pkg/util"
	"gorm.io/gorm"
)

type ActivityLogService struct {
	db *gorm.DB
}

func NewActivityLogService(db *gorm.DB) *ActivityLogService {
	return &ActivityLogService{db: db}
}

func (activityLogService *ActivityLogService) CreateActivityLog(activityLog *ActivityLog) error {
	return activityLogService.db.Create(activityLog).Error
}

func (activityLogService *ActivityLogService) GetActivityLogs() (*[]*ActivityLog, error) {
	var activityLogs []*ActivityLog
	err := activityLogService.db.Find(&activityLogs).Error
	if err != nil {
		return nil, err
	}
	return &activityLogs, nil
}

func CreateActivityLogTest(db *gorm.DB, message string) error {
	activityLogService := NewActivityLogService(db)
	activityLog := ActivityLog{Uuid: util.GenerateUUID(), Category: "test", Message: message}
	err := activityLogService.CreateActivityLog(&activityLog)
	if err != nil {
		return nil
	}
	return err
}
