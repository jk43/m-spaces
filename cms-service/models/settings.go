package models

import (
	"errors"
	"time"

	"github.com/moly-space/molylibs/service"
	"gorm.io/gorm"
)

// ValueType enum type definition

// Settings struct - uses composite primary key
type Setting struct {
	ID        uint                            `json:"-" gorm:"primaryKey;autoIncrement"`
	BoardID   uint                            `json:"-" gorm:"primaryKey;not null;index;uniqueIndex:idx_board_k"`
	K         string                          `json:"key" gorm:"size:255;not null;uniqueIndex:idx_board_k"`
	ValueType service.ServiceSettingValueType `json:"v_type" gorm:"type:enum('string','int','float','boolean');default:'string'"`
	V         string                          `json:"v" gorm:"size:255;not null"`
	CreatedAt time.Time                       `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time                       `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt                  `json:"-" gorm:"index"`
}

func (repo *DBRepo) InsertSettings(s []Setting) error {
	for _, v := range s {
		if err := repo.Mysql.Create(&v).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repo *DBRepo) UpsertSettings(id uint, elems service.FormElems) error {
	for k, v := range elems {
		setting := Setting{
			BoardID:   id,
			K:         k,
			V:         v.Value,
			ValueType: v.ValueType,
		}

		// Check if existing record exists
		var existingSetting Setting
		err := repo.Mysql.Where("board_id = ? AND k = ?", id, k).First(&existingSetting).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new record if it doesn't exist
				repo.Mysql.Create(&setting)
			} else {
				return err
			}
		} else {
			// Update if record exists
			existingSetting.V = v.Value
			existingSetting.ValueType = v.ValueType
			repo.Mysql.Save(&existingSetting)
		}
	}
	return nil
}

func (repo *DBRepo) GetSettingsWithBoardIDKye(boardID uint, k string) (*Setting, error) {
	var settings Setting
	err := repo.Mysql.Where("board_id = ? AND k = ?", boardID, k).Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}
