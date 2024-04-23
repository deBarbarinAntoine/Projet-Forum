package models

import (
	"database/sql"
	"time"
)

type DbData interface {
	Create()
	Fetch(any)
	GetId(any) int
	Exists(any) bool
	Update(any)
}

var UserFields = struct {
	Id         string
	Username   string
	Email      string
	HashedPwd  string
	Salt       string
	AvatarPath string
	Role       string
	BirthDate  string
	CreatedAt  string
	UpdatedAt  string
	VisitedAt  string
	Bio        string
	Signature  string
	Status     string
}{
	Id:         "Id_users",
	Username:   "Username",
	Email:      "Email",
	HashedPwd:  "Password",
	Salt:       "Salt",
	AvatarPath: "Avatar_path",
	Role:       "Role",
	BirthDate:  "Birth_date",
	CreatedAt:  "Created_at",
	UpdatedAt:  "Updated_at",
	VisitedAt:  "Visited_at",
	Bio:        "Bio",
	Signature:  "Signature",
	Status:     "Status",
}

type User struct {
	Id         int            `json:"id_users"`
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	HashedPwd  sql.NullString `json:"password"`
	Salt       sql.NullString `json:"salt"`
	AvatarPath sql.NullString `json:"avatar_path"`
	Role       string         `json:"role"`
	BirthDate  sql.NullTime   `json:"birth_date"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	VisitedAt  sql.NullTime   `json:"visited_at"`
	Bio        sql.NullString `json:"bio"`
	Signature  sql.NullString `json:"signature"`
	Status     sql.NullString `json:"status"`
}

type Category struct {
	Id               int       `json:"id_categories"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IdAuthor         int       `json:"id_author"`
	IdParentCategory int       `json:"id_parent_category"` // fixme handle null values
}

type Tag struct {
	Id        int       `json:"id_tags"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IdAuthor  int       `json:"id_author"`
}

type Thread struct {
	Id          int       `json:"id_threads"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	IdAuthor    int       `json:"id_author"`
	IdCategory  int       `json:"id_categories"`
}

type Post struct {
	Id           int       `json:"id_posts"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IdAuthor     int       `json:"id_author"`
	IdParentPost int       `json:"id_parent_post"`
	IdThread     int       `json:"id_thread"`
}

type Friend struct {
	Id        int       `json:"id_friends"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
}
