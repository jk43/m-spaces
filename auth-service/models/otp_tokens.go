package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type OTPToken struct {
	Token  sql.NullString
	UserID string
	Code   string `json:"code" validate:"required"`
	gorm.Model
}

func (repo *DBRepo) InsertOTPToken(token *OTPToken) (uint, error) {
	res := repo.Mysql.Create(token)
	if res.Error != nil {
		return 0, res.Error
	}
	return token.ID, nil
}

func (repo *DBRepo) GetOTPTokenWithToken(token string) (*OTPToken, error) {
	otpToken := OTPToken{}
	res := repo.Mysql.First(&otpToken, "token = ?", token)
	if res.Error != nil {
		return nil, res.Error
	}
	return &otpToken, nil
}

func (repo *DBRepo) DeleteOTPToken(token string) error {
	res := repo.Mysql.Unscoped().Delete(&OTPToken{}, "token = ?", token)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) DeleteOTPTokenWithInterval(interval int) error {
	res := repo.Mysql.Unscoped().Delete(&OTPToken{}, "created_at < ?", time.Now().Add(-time.Duration(interval)*time.Minute))
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) GetOTPTokenWithTokenAndCode(token, code string) (*OTPToken, error) {
	otpToken := OTPToken{}
	res := repo.Mysql.First(&otpToken, "token = ? AND code = ?", token, code)
	if res.Error != nil {
		return nil, res.Error
	}
	return &otpToken, nil
}
