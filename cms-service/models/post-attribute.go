package models

import "time"

type PostAttribute struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;type:bigint unsigned"`
	PostID    uint64    `json:"post_id" gorm:"type:bigint unsigned;not null;index:post_id"`
	Key       string    `json:"key" gorm:"type:varchar(100);not null;index:key"`
	Value     string    `json:"value" gorm:"type:text"`
	ValueType string    `json:"value_type" gorm:"type:varchar(20);default:'string'"` // string, number, boolean, json
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(3);not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime(3)"`

	Post Post `json:"post" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
