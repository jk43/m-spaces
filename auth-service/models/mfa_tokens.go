package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type MFAToken struct {
	Token  sql.NullString
	UserID string
	Code   string `json:"code" validate:"required"`
	gorm.Model
}

func (repo *DBRepo) InsertMFAToken(token *MFAToken) (uint, error) {
	res := repo.Mysql.Create(token)
	if res.Error != nil {
		return 0, res.Error
	}
	return token.ID, nil
}

func (repo *DBRepo) GetMFATokenWithToken(token string) (*MFAToken, error) {
	mfaToken := MFAToken{}
	res := repo.Mysql.First(&mfaToken, "token = ?", token)
	if res.Error != nil {
		return nil, res.Error
	}
	return &mfaToken, nil
}

func (repo *DBRepo) DeleteMFAToken(token string) error {
	res := repo.Mysql.Unscoped().Delete(&MFAToken{}, "token = ?", token)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) DeleteMFATokenWithInterval(interval int) error {
	res := repo.Mysql.Unscoped().Delete(&MFAToken{}, "created_at < ?", time.Now().Add(-time.Duration(interval)*time.Minute))
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) GetMFATokenWithTokenAndCode(token, code string) (*MFAToken, error) {
	mfaToken := MFAToken{}
	res := repo.Mysql.First(&mfaToken, "token = ? AND code = ?", token, code)
	if res.Error != nil {
		return nil, res.Error
	}
	return &mfaToken, nil
}
