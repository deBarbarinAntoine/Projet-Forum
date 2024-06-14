package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
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
