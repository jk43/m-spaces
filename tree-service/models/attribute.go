package models

import (
	"fmt"
	"time"

	"github.com/moly-space/molylibs/utils"
)

// TreeAttribute represents the structure of the `tree_attributes` table, including foreign key constraints and JSON tags.
type TreeAttribute struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;<:UNSIGNED>" json:"id"`
	OrganizationID string    `gorm:"size:255;not null;index:,type:btree" json:"organization_id"`
	Slug           string    `gorm:"unique;size:255;not null" json:"slug"`
	TreeID         uint      `gorm:"not null;index:idx_tree_id,type:btree;<:UNSIGNED>;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tree_id"`
	Label          string    `gorm:"size:255;not null" json:"label"`
	View           string    `gorm:"size:255;not null" json:"view"`
	Description    string    `gorm:"type:text" json:"description"`
	Options        *string   `gorm:"type:json" json:"options"` // Assuming JSON support in the database
	CreatedBy      string    `gorm:"size:255;not null" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type TreeOptions struct {
}

func (repo *DBRepo) GetTreeAttributeWithSlug(slug string) (*TreeAttribute, error) {
	attr := &TreeAttribute{
		Slug: slug,
	}
	res := repo.Mysql.First(attr, "slug = ?", slug)

	if res.Error != nil {
		return nil, res.Error
	}
	return attr, nil
}

func (repo *DBRepo) GetTreeAttributeWithOrgIDAndSlug(orgID, slug string) (*TreeAttribute, error) {
	attr := &TreeAttribute{
		Slug: slug,
	}
	res := repo.Mysql.First(attr, "slug = ? AND organization_id = ?", slug, orgID)

	if res.Error != nil {
		return nil, res.Error
	}
	return attr, nil
}

func (repo *DBRepo) InsertTreeAttribute(attr *TreeAttribute) (uint, error) {
	res := repo.Mysql.Create(attr)

	if res.Error != nil {
		return 0, res.Error
	}
	return attr.ID, nil
}

func (repo *DBRepo) UpdateTreeAttribute(orgID string, attr *TreeAttribute) error {

	res := repo.Mysql.Model(&TreeAttribute{}).Where("organization_id = ? AND id = ?", orgID, attr.ID).Updates(attr)
	utils.TermDebugging(`res`, res)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("0 row affected")
	}
	return nil
}
