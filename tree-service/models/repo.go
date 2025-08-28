package models

import "gorm.io/gorm"

type DBRepo struct {
	Mysql *gorm.DB
}

type DatabaseRepo interface {
	GetTreeAttributeWithSlug(slug string) (*TreeAttribute, error)
	GetTreeAttributeWithOrgIDAndSlug(orgID, slug string) (*TreeAttribute, error)
	InsertTreeAttribute(attr *TreeAttribute) (uint, error)
	InsertTree(tree *Tree) (*Tree, error)
	UpdateTreeWithID(id uint, tree *Tree) error
	DeleteTree(t *Tree) error
	DeleteTreeWithID(id uint) error
	GetRootsWithOrgID(orgID string) ([]*Tree, error)
	GetTreeWithSlug(orgID string, slug string) (*Tree, error)
	GetTreeWithID(orgID string, id uint) (*Tree, error)
	GetChildrenWithSlug(orgID string, slug string) ([]*Tree, error)
	UpdateOrder(t *Tree, i uint) error
	UpdateTreeAttribute(orgID string, attr *TreeAttribute) error
}
