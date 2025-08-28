package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	BoardID   string
	ParentID  *uint `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ParentID;references:ID" json:"parent_id"`
	Message   string
	CreatedBy string `gorm:"index:idx_createdby"`
	UpdatedAt time.Time
	CreatedAt time.Time
	gorm.Model
}
