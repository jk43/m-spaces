package models

import (
	"database/sql"

	"github.com/moly-space/molylibs/utils"
	"gorm.io/gorm"
)

// type YesOrNo string

// const (
// 	Yes YesOrNo = "Y"
// 	No  YesOrNo = "N"
// )

type VerificationTokens struct {
	Token         sql.NullString
	Verified      utils.YesOrNo `gorm:"type:ENUM('Y', 'N');default:'N'"`
	ResetPassword utils.YesOrNo `gorm:"type:ENUM('Y', 'N');default:'N'"`
	CredentialsID uint
	gorm.Model
}

func (repo *DBRepo) GetVerificationTokenWithToken(token string) (*VerificationTokens, error) {
	vToken := VerificationTokens{}
	res := repo.Mysql.First(&vToken, "token = ?", token)
	if res.Error != nil {
		return nil, res.Error
	}
	return &vToken, nil
}

func (repo *DBRepo) UpdateVerificationToken(token string, vt *VerificationTokens) error {
	return repo.Mysql.Model(&VerificationTokens{}).Where("token = ?", token).Updates(vt).Error
}
