package models

import "gorm.io/gorm"

type DBRepo struct {
	Mysql *gorm.DB
}

type DatabaseRepo interface {
	InsertCredentials(cred *Credentials) (uint, error)
	GetCredentialsWithEmailAndOrgID(email string, orgID string) (*Credentials, error)
	GetCredentialsWithUserID(id string) (*Credentials, error)
	UpdateCredentialsWithCredentials(credentials *Credentials) error
	GetVerificationTokenWithToken(token string) (*VerificationTokens, error)
	UpdateVerificationToken(token string, vt *VerificationTokens) error
	FindWithEmailAndOrgID(email, orgID string) (*Credentials, error)
	UpdatePasswordWithUserID(id, password, salt string) error
	UpdateEmailWithUserID(id string, email string) error
	UpdateNameWithUserID(id string, firstName string, lastName string) error
	InsertResetPasswordToken(rpt *ResetPasswordTokens) (uint, error)
	GetResetPasswordTokenWithToken(token string) (*ResetPasswordTokens, error)
	GetCredentialsWithVeriToken(token string) (*Credentials, error)
	CleanUpResetPassword(minutes int) error
	DeleteResetPasswordToken(id uint) error
	DeleteCredentialsWithUserID(id string) error
	DeleteResetPasswordTokenWithUserID(id string) error
	InsertMFAToken(token *MFAToken) (uint, error)
	GetMFATokenWithToken(token string) (*MFAToken, error)
	DeleteMFAToken(token string) error
	DeleteMFATokenWithInterval(interval int) error
	GetMFATokenWithTokenAndCode(token, code string) (*MFAToken, error)
	InsertOTPToken(token *OTPToken) (uint, error)
	GetOTPTokenWithToken(token string) (*OTPToken, error)
	DeleteOTPToken(token string) error
	DeleteOTPTokenWithInterval(interval int) error
	GetOTPTokenWithTokenAndCode(token, code string) (*OTPToken, error)
}
