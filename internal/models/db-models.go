package models

import (
	"database/sql"
	"time"
)

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
