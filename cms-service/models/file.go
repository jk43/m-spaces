package models

import (
	"time"

	"github.com/moly-space/molylibs/utils"
	"gorm.io/gorm"
)

type File struct {
	ID               uint64          `gorm:"type:bigint;primaryKey;autoIncrement" json:"-"`
	UserID           string          `gorm:"type:varchar(255);not null;index" json:"-"`
	PostID           *uint64         `gorm:"type:bigint;index" json:"-"`
	FileID           string          `gorm:"type:varchar(255);not null;index" json:"id"`
	FormID           string          `gorm:"type:varchar(255);index" json:"formId"`
	ContentType      string          `gorm:"type:varchar(255)" json:"contentType"`
	OriginalFileName string          `gorm:"type:varchar(255)" json:"name"`
	S3Path           string          `gorm:"type:varchar(255)" json:"s3Path"`
	Service          string          `gorm:"type:varchar(255)" json:"service"`
	ServiceCtx       string          `gorm:"type:varchar(255)" json:"serviceCtx"`
	Size             uint64          `gorm:"type:bigint" json:"size"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"-"`
	DeletedAt        *gorm.DeletedAt `gorm:"index" json:"-"`
}

func (repo *DBRepo) InsertFile(f *File) (uint64, error) {
	res := repo.Mysql.Create(f)
	if res.Error != nil {
		return 0, res.Error
	}
	return f.ID, nil
}

func (repo *DBRepo) GetOrphanFilesWithUserID(userID string) ([]*File, error) {
	var files []*File
	res := repo.Mysql.Where("user_id = ? AND post_id IS NULL", userID).Find(&files)
	if res.Error != nil {
		return nil, res.Error
	}
	return files, nil
}

func (repo *DBRepo) GetFilesWithPostID(id uint64) ([]*File, error) {
	var files []*File
	res := repo.Mysql.Where("post_id = ?", id).Find(&files)
	if res.Error != nil {
		return nil, res.Error
	}
	return files, nil
}

func (repo *DBRepo) UpdateFilePostIDWithFormID(formID string, postID uint64) error {
	res := repo.Mysql.Model(&File{}).Where("form_id = ?", formID).Update("post_id", postID)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Update form id for orphan files
func (repo *DBRepo) UpdateFileFormIDWithUserID(userID string, formID string) error {
	utils.TermDebugging(`req.Payload.Data["formId"][0]`, formID)
	res := repo.Mysql.Model(&File{}).Where("user_id = ? AND post_id IS NULL", userID).Update("form_id", formID)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) DeleteFile(f *File) error {
	res := repo.Mysql.Delete(f)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) GetFile(f *File) error {
	res := repo.Mysql.Where("file_id = ? AND user_id = ?", f.FileID, f.UserID).Find(f)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
