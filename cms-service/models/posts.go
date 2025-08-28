package models

import (
	"time"

	"github.com/moly-space/molylibs/utils"
	"gorm.io/gorm"
)

type UserName struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Method string

const (
	MethodPost   Method = "POST"
	MethodPut    Method = "PUT"
	MethodDelete Method = "DELETE"
)

// Post struct - corresponds to posts table
type Post struct {
	ID        uint64          `json:"-" gorm:"primaryKey;autoIncrement;type:bigint unsigned"`
	BoardID   uint            `json:"-" gorm:"primaryKey;type:int unsigned;not null;index:board_id"`
	RootID    *uint64         `json:"-" gorm:"type:bigint unsigned;index:root_id"`
	ParentID  *uint64         `json:"-" gorm:"type:bigint unsigned;index:parent_id;foreignKey:ID;constraint:OnDelete:CASCADE"`
	Slug      string          `json:"slug" gorm:"type:varchar(255);not null;uniqueIndex:slug"`
	Title     string          `json:"title" gorm:"type:varchar(255);index:title"`
	Text      string          `json:"text" gorm:"type:text;fulltext:text"`
	CreatedAt time.Time       `json:"created_at" gorm:"type:datetime(3);not null"`
	UpdatedAt time.Time       `json:"-" gorm:"type:datetime(3)"`
	DeletedAt *gorm.DeletedAt `json:"-" gorm:"type:datetime(3)"`
	// Relationship settings
	Board          Board            `json:"-" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Parent         *Post            `json:"-" gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE"`
	Children       []*Post          `json:"-" gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE"`
	UserName       *UserName        `json:"user_name" gorm:"-"`
	TotalComments  int              `json:"total_comments" gorm:"-"`
	Files          []*File          `json:"-"`
	PostAuthor     *PostAuthor      `json:"-" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	PostAttributes []*PostAttribute `json:"-" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

func (repo *DBRepo) InsertPost(p *Post) (uint64, error) {
	res := repo.Mysql.Create(p)
	if res.Error != nil {
		return 0, res.Error
	}
	return p.ID, nil
}

func (repo *DBRepo) UpdatePostParent(p *Post) error {
	res := repo.Mysql.Model(&Post{}).Where("id = ?", p.ID).Updates(p)
	return res.Error
}

func (repo *DBRepo) GetPostsWithBoardIDAndPagination(boardID uint, pagination *utils.Pagination, search string, sortBy string, descending bool) ([]*Post, int64, error) {
	var posts []*Post
	var total int64

	query := repo.Mysql.Model(&Post{}).Omit("text").Where("board_id = ? AND id = parent_id", boardID)

	// Add search conditions
	if search != "" {
		query = query.Where("title LIKE ? OR text LIKE ?", "%"+search+"%", "%"+search+"%")
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
	if err := query.Preload("PostAuthor").Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (repo *DBRepo) UpdatePost(p *Post) error {
	return repo.Mysql.Model(&Post{}).Where("id = ?", p.ID).Updates(p).Error
}

func (repo *DBRepo) GetPost(p *Post) error {
	post := repo.Mysql.Model(&Post{})
	if p.BoardID != 0 {
		post = post.Where("board_id = ?", p.BoardID)
	}
	if p.Slug != "" {
		post = post.Where("slug = ?", p.Slug)
	}
	return post.First(p).Error
}

func (repo *DBRepo) DeletePost(p *Post) error {
	return repo.Mysql.Model(&Post{}).Where("slug = ?", p.Slug).Delete(p).Error
}

func (repo *DBRepo) GetPostSlugWithSlugAndParams(slug string, params ...string) (string, error) {
	post := Post{}

	query := repo.Mysql.Where("slug = ?", slug)
	err := query.First(&post).Error
	return post.Slug, err
}

func (repo *DBRepo) GetTotalCommentsWithID(id uint64) (int, error) {
	var total int64

	err := repo.Mysql.Model(&Post{}).Where("root_id = ?", id).Count(&total).Error
	return int(total - 1), err
}

func (repo *DBRepo) GetCommentsWithIDAndPagination(id uint64, pagination *utils.Pagination, search string, sortBy string, descending bool) ([]*Post, int64, error) {
	var posts []*Post
	var total int64

	query := repo.Mysql.Model(&Post{}).Where("parent_id = ? AND id != parent_id", id)

	// 검색 조건 추가
	if search != "" {
		query = query.Where("title LIKE ? OR text LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 전체 개수 조회
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 정렬 조건 추가
	if sortBy != "" {
		orderClause := sortBy
		if descending {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}
		query = query.Order(orderClause)
	}

	// 페이지네이션 적용
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.RowsPerPage
		query = query.Offset(int(offset)).Limit(int(pagination.RowsPerPage))
	}

	// 데이터 조회
	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
