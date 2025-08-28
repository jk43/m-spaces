package models

import (
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"gorm.io/gorm"
)

type DBRepo struct {
	Mysql *gorm.DB
}

type DatabaseRepo interface {
	InsertBoard(b *Board) (uint, error)
	GetBoardSlugWithSlug(slug string, params ...string) (string, error)
	GetBoardsWithOrgID(orgID string, pagination *utils.Pagination, search string) ([]Board, error)
	GetBoardsWithOrgIDAndPagination(orgID string, pagination *utils.Pagination, search string, sortBy string, descending bool) ([]Board, int64, error)
	UpsertSettings(id uint, elems service.FormElems) error
	GetSettingsWithBoardIDKye(boardID uint, k string) (*Setting, error)
	UpdateBoard(b *Board) error
	DeleteBoard(b *Board) error
	GetBoard(b *Board) error
	InsertPost(p *Post) (uint64, error)
	UpdatePostParent(p *Post) error
	GetPostsWithBoardIDAndPagination(boardID uint, pagination *utils.Pagination, search string, sortBy string, descending bool) ([]*Post, int64, error)
	UpdatePost(p *Post) error
	GetPost(p *Post) error
	DeletePost(p *Post) error
	GetPostSlugWithSlugAndParams(slug string, params ...string) (string, error)
	GetTotalCommentsWithID(id uint64) (int, error)
	GetCommentsWithIDAndPagination(id uint64, pagination *utils.Pagination, search string, sortBy string, descending bool) ([]*Post, int64, error)
	InsertFile(f *File) (uint64, error)
	GetFilesWithPostID(id uint64) ([]*File, error)
	UpdateFilePostIDWithFormID(formID string, postID uint64) error
	DeleteFile(f *File) error
	GetOrphanFilesWithUserID(userID string) ([]*File, error)
	UpdateFileFormIDWithUserID(userID string, formID string) error
	GetFile(f *File) error
	InsertPostAuthor(p *PostAuthor) error
	InsertGuestVerificationCode(g *GuestVerificationCode) error
	GetGuestVerificationCode(g *GuestVerificationCode) error
	DeleteGuestVerificationCode() error
}
