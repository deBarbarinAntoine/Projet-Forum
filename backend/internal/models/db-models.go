package models

import (
	"database/sql"
	"time"
)

// DbData provides all methods to interact with the database's items
type DbData interface {
	Create() int
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
	HashedPwd  sql.NullString `json:"password"` // fixme
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
	Id             int       `json:"id_categories"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Author         string    `json:"author"`
	ParentCategory string    `json:"parent_category"` // fixme handle null values
}

type Tag struct {
	Id        int       `json:"id_tags"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    string    `json:"author"`
}

type Thread struct {
	Id          int       `json:"id_threads"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	Author      string    `json:"author"`
	Category    string    `json:"categories"`
}

type Post struct {
	Id           int       `json:"id_posts"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Author       string    `json:"author"`
	ParentPostId int       `json:"parent_post_id"`
	ThreadId     int       `json:"thread_id"`
}

type Friend struct {
	Id        int       `json:"id_friends"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
}
