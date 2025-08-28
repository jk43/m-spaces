package models

type Tree struct {
	ID       uint  `gorm:"primaryKey;<:UNSIGNED>" json:"id"`
	RootID   uint  `gorm:"not null;index:idx_root_id,type:btree;<:UNSIGNED>;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"root_id"`
	ParentID *uint `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ParentID;references:ID" json:"parent_id"`
	Order    uint  `gorm:"not null;<:UNSIGNED>" json:"order"`
	// Relationships
	Attributes TreeAttribute `gorm:"foreignKey:TreeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"attributes"`
	Children   []*Tree       `gorm:"-" json:"children"`
}

func (repo *DBRepo) InsertTree(tree *Tree) (*Tree, error) {
	res := repo.Mysql.Create(tree)

	if res.Error != nil {
		return nil, res.Error
	}
	return tree, nil
}

func (repo *DBRepo) UpdateTreeWithID(id uint, tree *Tree) error {
	res := repo.Mysql.Model(&Tree{}).Where("id = ?", id).Updates(tree)

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *DBRepo) DeleteTreeWithID(id uint) error {
	err := repo.Mysql.Delete(&Tree{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepo) DeleteTree(t *Tree) error {
	err := repo.Mysql.Delete(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepo) GetRootsWithOrgID(orgID string) ([]*Tree, error) {
	var trees []*Tree
	// err := repo.Mysql.Joins("JOIN tree_attributes ON trees.id = tree_attributes.tree_id").
	// 	Where("tree_attributes.organization_id = ? AND trees.parent_id IS NULL", orgID).
	// 	Find(&trees).Error
	err := repo.Mysql.Preload("Attributes").Model(&Tree{}).Joins("JOIN tree_attributes ON trees.id = tree_attributes.tree_id").
		Where("tree_attributes.organization_id = ? AND trees.parent_id IS NULL", orgID).
		Find(&trees).Error
	if err != nil {
		return nil, err
	}
	return trees, nil
}

func (repo *DBRepo) GetTreeWithSlug(orgID string, slug string) (*Tree, error) {
	var tree Tree
	err := repo.Mysql.Preload("Attributes").Model(&Tree{}).Joins("JOIN tree_attributes ON trees.id = tree_attributes.tree_id").
		Where("tree_attributes.organization_id = ? AND tree_attributes.slug = ?", orgID, slug).
		First(&tree).Error
	if err != nil {
		return nil, err
	}
	return &tree, nil
}

func (repo *DBRepo) GetTreeWithID(orgID string, id uint) (*Tree, error) {
	var tree Tree
	err := repo.Mysql.Preload("Attributes").Model(&Tree{}).Joins("JOIN tree_attributes ON trees.id = tree_attributes.tree_id").
		Where("tree_attributes.organization_id = ? AND trees.id = ?", orgID, id).
		First(&tree).Error
	if err != nil {
		return nil, err
	}
	return &tree, nil
}

func (repo *DBRepo) GetChildrenWithSlug(orgID string, slug string) ([]*Tree, error) {
	var trees []*Tree
	err := repo.Mysql.Preload("Attributes").
		Model(&Tree{}).
		Select("t.*").
		Joins("JOIN tree_attributes ON trees.id = tree_attributes.tree_id").
		Joins("JOIN trees t on trees.id = t.parent_id").
		Joins("JOIN tree_attributes ta on t.id = ta.tree_id").
		Where("tree_attributes.organization_id = ? AND tree_attributes.slug = ?", orgID, slug).
		Order("t.order ASC").
		Find(&trees).
		Error
	if err != nil {
		return nil, err
	}
	return trees, nil
}

func (repo *DBRepo) UpdateOrder(t *Tree, i uint) error {
	err := repo.Mysql.Model(t).Update("order", i).Error
	if err != nil {
		return err
	}
	return nil
}

// func GetChildren(orgID string, id uint) ([]*Tree, error) {
// 	// Get the children of the given node
// 	// and return the list of children
// 	return nil, nil
// }

// func GetParent(orgID string, id uint) (*Tree, error) {
// 	// Get the parent of the given node
// 	// and return the parent node
// 	return nil, nil
// }

// func GetAncestors(orgID string, id uint) ([]*Tree, error) {
// 	// Get the ancestors of the given node
// 	// and return the list of ancestors
// 	return nil, nil
// }

// func GetSibling(orgID string, id uint) ([]*Tree, error) {
// 	// Get the siblings of the given node
// 	// and return the list of siblings
// 	return nil, nil
// }

// func DeleteTree(orgID string, id uint) error {
// 	// Delete the tree node
// 	return nil
// }
