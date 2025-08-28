package models

import "gorm.io/gorm"

type DBRepo struct {
	Mysql *gorm.DB
}

type DatabaseRepo interface {
	InsertFile(file *FileInfo) (uint64, error)
	GetFile(file *FileInfo) error
	DeleteFile(file *FileInfo) error
}
