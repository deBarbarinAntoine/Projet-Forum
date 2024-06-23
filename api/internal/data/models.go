package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrEditConflict      = errors.New("edit conflict")
	ErrDuplicateEmail    = errors.New("duplicate user email")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateFriend   = errors.New("duplicate friend")
	ErrDuplicateName     = errors.New("duplicate tag or category name")
	ErrDuplicateTitle    = errors.New("duplicate thread title")
	ErrDuplicateToken    = errors.New("duplicate token")
	ErrDuplicateEntry    = errors.New("duplicate entry")
)

type Models struct {
	Categories  CategoryModel
	Threads     ThreadModel
	Tags        TagModel
	Posts       PostModel
	Permissions PermissionModel
	Tokens      TokenModel
	Users       UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Categories:  CategoryModel{DB: db},
		Threads:     ThreadModel{DB: db},
		Tags:        TagModel{DB: db},
		Posts:       PostModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
