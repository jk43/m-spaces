package models

import (
	"time"

	"gorm.io/gorm"
)

type Contents struct {
	BoardID   string
	Message   string
	CreatedBy string `gorm:"index:idx_createdby"`
	UpdatedAt time.Time
	CreatedAt time.Time
	gorm.Model
}
