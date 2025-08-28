package models

import (
	"os"
	"time"

	"gorm.io/gorm"
)

type FileInfo struct {
	ID               uint64          `json:"-" gorm:"primaryKey;autoIncrement;type:bigint unsigned"`
	UUID             string          `json:"uuid" gorm:"index:idx_uuid,unique;not null"`
	OSFile           *os.File        `json:"-" gorm:"-"`
	BatchID          string          `json:"-" gorm:"index:idx_batchid"`
	OriginalFileName string          `json:"name"`
	FileName         string          `gorm:"index:idx_filename,unique" json:"fileName"` //hashed file name
	S3Key            string          `json:"s3Path"`
	ContentType      string          `json:"contentType"`
	Size             int64           `json:"size"`
	Service          string          `json:"service"`
	ServiceCtx       string          `json:"serviceCtx"`
	FormID           string          `json:"-" gorm:"index:idx_formid"` //Id from the form to send service. Optional
	FormKey          string          `json:"-"`                         //Key name of file
	UserID           string          `json:"-" gorm:"index:idx_userid" json:"UserId"`
	OrganizationID   string          `json:"-" gorm:"index:idx_organizationid" json:"OrganizationId"`
	IP               string          `json:"-"`
	Error            bool            `json:"-"`
	ErrorMessage     string          `json:"-"`
	CreatedAt        time.Time       `json:"created_at" gorm:"type:datetime(3);not null"`
	UpdatedAt        time.Time       `json:"-" gorm:"type:datetime(3)"`
	DeletedAt        *gorm.DeletedAt `json:"-" gorm:"type:datetime(3)"`
}

func (repo *DBRepo) InsertFile(file *FileInfo) (uint64, error) {
	res := repo.Mysql.Create(file)
	if res.Error != nil {
		return 0, res.Error
	}
	return file.ID, nil
}

func (repo *DBRepo) GetFile(file *FileInfo) error {
	res := repo.Mysql.Where("uuid = ? AND user_id = ? AND organization_id = ?", file.UUID, file.UserID, file.OrganizationID).Find(file)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) DeleteFile(file *FileInfo) error {
	res := repo.Mysql.Where("uuid = ? AND user_id = ? AND organization_id = ?", file.UUID, file.UserID, file.OrganizationID).Delete(file)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
