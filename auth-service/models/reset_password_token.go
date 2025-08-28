package models

import (
	"strconv"

	"gorm.io/gorm"
)

type ResetPasswordTokens struct {
	UserID            string
	VerificationToken string `gorm:"index:idx_verification_token,unique"`
	gorm.Model
}

func (repo *DBRepo) InsertResetPasswordToken(rpt *ResetPasswordTokens) (uint, error) {
	res := repo.Mysql.Create(rpt)
	if res.Error != nil {
		return 0, res.Error
	}
	return rpt.ID, nil
}

func (repo *DBRepo) GetResetPasswordTokenWithToken(token string) (*ResetPasswordTokens, error) {
	var rpt ResetPasswordTokens
	res := repo.Mysql.Where("verification_token = ?", token).First(&rpt)
	if res.Error != nil {
		return nil, res.Error
	}
	return &rpt, nil
}

func (repo *DBRepo) DeleteResetPasswordToken(id uint) error {
	res := repo.Mysql.Delete(&ResetPasswordTokens{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) CleanUpResetPassword(minutes int) error {
	res := repo.Mysql.Where("created_at < DATE_SUB(NOW(), INTERVAL " + strconv.Itoa(minutes) + " MINUTE)").Delete(&ResetPasswordTokens{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) DeleteResetPasswordTokenWithUserID(id string) error {
	return repo.Mysql.Delete(&ResetPasswordTokens{}, "user_id = ?", id).Error
}
