package models

import "gorm.io/gorm"

type DBRepo struct {
	Mysql *gorm.DB
}

type DatabaseRepo interface {
}
