package quoteApp

import "gorm.io/gorm"

type ActivityLog struct {
	gorm.Model
	Uuid     string `json:"uuid"`
	Category string `json:"category"`
	Message  string `json:"message"`
}
