package models

import (
	"time"

	"gorm.io/gorm"
)

type PostAuthor struct {
	ID        uint64          `json:"id" gorm:"primaryKey;autoIncrement;type:bigint unsigned"`
	PostID    uint64          `json:"post_id" gorm:"type:bigint unsigned;not null;index:post_id"`
	UserID    string          `json:"user_id" gorm:"type:varchar(32);index:user_id"`
	Name      string          `json:"name" gorm:"type:varchar(255)"`
	Email     string          `json:"email" gorm:"type:varchar(255);"`
	CreatedAt time.Time       `json:"created_at" gorm:"type:datetime(3);not null"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:datetime(3)"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime(3)"`

	Post Post `json:"post" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

func (repo *DBRepo) InsertPostAuthor(p *PostAuthor) error {
	res := repo.Mysql.Create(p)
	return res.Error
}
