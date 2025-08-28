package models

import (
	"time"

	"github.com/moly-space/molylibs/utils"
)

type GuestVerificationCode struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Code      string    `json:"code" gorm:"type:varchar(255)"`
	Email     string    `json:"email" gorm:"type:varchar(255)"`
	PostID    *uint64   `json:"post_id" gorm:"type:bigint unsigned"`
	Method    Method    `json:"method" gorm:"type:ENUM('POST', 'PUT', 'DELETE')"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (repo *DBRepo) InsertGuestVerificationCode(g *GuestVerificationCode) error {
	repo.DeleteGuestVerificationCode()
	g.Code = utils.GetVerificationCode()
	return repo.Mysql.Create(g).Error
}

func (repo *DBRepo) GetGuestVerificationCode(g *GuestVerificationCode) error {
	repo.DeleteGuestVerificationCode()
	if g.PostID == nil {
		return repo.Mysql.Where("email = ? AND post_id is NULL AND method = ?", g.Email, g.Method).Last(g).Error
	}
	return repo.Mysql.Where("email = ? AND post_id = ? AND method = ?", g.Email, g.PostID, g.Method).Last(g).Error
}

func (repo *DBRepo) DeleteGuestVerificationCode() error {
	return repo.Mysql.Where("created_at < DATE_SUB(NOW(), INTERVAL 5 MINUTE)").Delete(&GuestVerificationCode{}).Error
}
