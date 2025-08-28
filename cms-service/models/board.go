package models

import (
	"time"

	"github.com/moly-space/molylibs/utils"
	"gorm.io/gorm"
)

// Board struct - corresponds to boards table
type Board struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned"`
	OrganizationID string         `json:"organization_id" gorm:"type:varchar(255);not null;index:idx_organizationid;uniqueIndex:idx_org_slug"`
	Name           string         `json:"name" gorm:"type:varchar(255);not null"`
	Slug           string         `json:"slug" gorm:"type:varchar(255);not null;uniqueIndex:idx_org_slug"`
	Active         utils.YesOrNo  `json:"active" gorm:"type:ENUM('Y', 'N');default:'N'"`
	Settings       []Setting      `json:"settings" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	CreatedBy      string         `json:"-" gorm:"type:varchar(255)"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	gorm.Model
}

func (b *Board) GetSetting(key, value string) string {
	for _, setting := range b.Settings {
		if setting.K == key {
			return setting.V
		}
	}
	return value
}

// for slug
func (repo *DBRepo) GetBoardSlugWithSlug(slug string, params ...string) (string, error) {
	board := Board{}

	query := repo.Mysql.Preload("Settings").Where("slug = ? AND organization_id = ?", slug, params[0])
	err := query.First(&board).Error
	return board.Slug, err
}

func (repo *DBRepo) GetBoard(b *Board) error {
	var query *gorm.DB
	if b.ID != 0 {
		query = repo.Mysql.Preload("Settings").Where("id = ? AND organization_id = ?", b.ID, b.OrganizationID)

	} else {
		query = repo.Mysql.Preload("Settings").Where("slug = ? AND organization_id = ?", b.Slug, b.OrganizationID)
	}
	err := query.First(b).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepo) InsertBoard(b *Board) (uint, error) {
	res := repo.Mysql.Create(b)
	if res.Error != nil {
		return 0, res.Error
	}
	return b.ID, nil
}

func (repo *DBRepo) GetBoardsWithOrgIDAndPagination(orgID string, pagination *utils.Pagination, search string, sortBy string, descending bool) ([]Board, int64, error) {
	var boards []Board
	var total int64

	query := repo.Mysql.Model(&Board{}).Where("organization_id = ?", orgID)

	// Add search conditions
	if search != "" {
		query = query.Where("name LIKE ? OR slug LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Add sorting conditions
	if sortBy != "" {
		orderClause := sortBy
		if descending {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}
		query = query.Order(orderClause)
	}

	// Apply pagination
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.RowsPerPage
		query = query.Offset(int(offset)).Limit(int(pagination.RowsPerPage))
	}

	// Fetch data
	if err := query.Preload("Settings").Find(&boards).Error; err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

// func (repo *DBRepo) GetBoardsByOrgIDWithPagination(orgID string, pagination *utils.Pagination, search string) ([]Board, int64, error) {
// 	var boards []Board
// 	var total int64

// 	query := repo.Mysql.Model(&Board{}).Where("organization_id = ?", orgID)

// 	// Add search conditions
// 	if search != "" {
// 		query = query.Where("name LIKE ? OR slug LIKE ?", "%"+search+"%", "%"+search+"%")
// 	}

// 	// Get total count
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	// Apply pagination
// 	if pagination != nil {
// 		offset := (pagination.Page - 1) * pagination.Limit
// 		query = query.Offset(offset).Limit(pagination.Limit)
// 	}

// 	// Fetch data
// 	if err := query.Preload("Settings").Find(&boards).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return boards, total, nil
// }

func (repo *DBRepo) GetBoardsWithOrgID(orgID string, pagination *utils.Pagination, search string) ([]Board, error) {
	var boards []Board
	res := repo.Mysql.Preload("Settings").Where("organization_id = ?", orgID).Find(&boards)
	if res.Error != nil {
		return nil, res.Error
	}
	return boards, nil
}

func (repo *DBRepo) UpdateBoard(b *Board) error {
	return repo.Mysql.Model(&Board{}).Where("id = ? AND organization_id = ?", b.ID, b.OrganizationID).Updates(b).Error
}

func (repo *DBRepo) DeleteBoard(b *Board) error {
	return repo.Mysql.Model(&Board{}).Where("id = ? AND organization_id = ?", b.ID, b.OrganizationID).Delete(b).Error
}
