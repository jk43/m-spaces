package models

import (
	"github.com/moly-space/molylibs/service"
	"gorm.io/gorm"
)

type Credentials struct {
	UserID             string
	Email              string             `gorm:"index:idx_account,unique" json:"email" validate:"required,email"`
	OrganizationID     string             `gorm:"index:idx_account,unique"`
	FirstName          string             `json:"firstName" validate:"required"`
	LastName           string             `json:"lastName" validate:"required"`
	Password           string             `json:"-" validate:"required"`
	Salt               string             `json:"-" validate:"required"`
	VerificationTokens VerificationTokens `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status             service.UserStatus `gorm:"type:ENUM('active', 'inactive', 'deleted', 'waiting');default:'inactive'"`
	gorm.Model
}

// type VerificationTokens struct {
// 	Token         sql.NullString
// 	Verified      YesOrNo `gorm:"type:ENUM('Y', 'N');default:'N'"`
// 	ResetPassword YesOrNo `gorm:"type:ENUM('Y', 'N');default:'N'"`
// 	CredentialsID uint
// 	gorm.Model
// }

func (repo *DBRepo) FindWithEmailAndOrgID(email, orgID string) (*Credentials, error) {
	cred := Credentials{}
	err := repo.Mysql.
		Preload("Credentials").
		First(&cred, "email = ? AND organization_id = ?", email, orgID).Error

	if err != nil {
		return nil, err
	}

	return &cred, nil
}

// func (repo *DBRepo) InsertCredentialsWithUser(user *User) (int, error) {
// 	res := repo.Mysql.Create(user)
// 	if res.Error != nil {
// 		return 0, res.Error
// 	}
// 	return int(user.ID), nil
// }

func (repo *DBRepo) InsertCredentials(cred *Credentials) (uint, error) {
	res := repo.Mysql.Create(cred)
	if res.Error != nil {
		return 0, res.Error
	}
	return cred.ID, nil
}

func (repo *DBRepo) GetCredentialsWithID(id string) (*Credentials, error) {
	cred := Credentials{}
	res := repo.Mysql.Preload("VerificationTokens").First(&cred, "id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &cred, nil
}

func (repo *DBRepo) GetCredentialsWithVeriToken(token string) (*Credentials, error) {
	cred := Credentials{}
	err := repo.Mysql.Preload("VerificationTokens").Table("credentials").Joins("join verification_tokens on verification_tokens.credentials_id = credentials.id").Where("verification_tokens.token = ?", token).First(&cred).Error
	//res := repo.Mysql.Preload("VerificationTokens").First(&cred, "verification_tokens.token = ?", token)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (repo *DBRepo) GetCredentialsWithEmailAndOrgID(email string, orgID string) (*Credentials, error) {
	creds := Credentials{}

	err := repo.Mysql.
		Preload("VerificationTokens").
		Where("credentials.email = ? AND credentials.organization_id = ?", email, orgID).
		First(&creds).Error

	if err != nil {
		return nil, err
	}
	return &creds, nil
}

func (repo *DBRepo) GetCredentialsWithUserID(id string) (*Credentials, error) {
	cred := Credentials{}
	res := repo.Mysql.Preload("VerificationTokens").First(&cred, "user_id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &cred, nil
}

func (repo *DBRepo) UpdateCredentialsWithCredentials(credentials *Credentials) error {
	repo.Mysql.Save(credentials)
	return nil
}

func (repo *DBRepo) UpdatePasswordWithUserID(id, password, salt string) error {
	return repo.Mysql.Model(&Credentials{}).Where("user_id = ?", id).Updates(Credentials{Password: password, Salt: salt}).Error
}

func (repo *DBRepo) UpdateEmailWithUserID(id string, email string) error {
	return repo.Mysql.Model(&Credentials{}).Where("user_id = ?", id).Update("email", email).Error
}

func (repo *DBRepo) UpdateNameWithUserID(id string, firstName string, lastName string) error {
	return repo.Mysql.Model(&Credentials{}).Where("user_id = ?", id).Updates(Credentials{FirstName: firstName, LastName: lastName}).Error
}

func (repo *DBRepo) DeleteCredentialsWithUserID(id string) error {
	return repo.Mysql.Unscoped().Delete(&Credentials{}, "user_id = ?", id).Error
}
