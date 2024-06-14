package models

import "Projet-Forum/internal/db"

// User methods

func (user *User) GetSqlRows() []any {
	return []any{&user.Id, &user.Username, &user.Email, &user.HashedPwd, &user.Salt, &user.AvatarPath, &user.Role, &user.BirthDate, &user.CreatedAt, &user.UpdatedAt, &user.VisitedAt, &user.Bio, &user.Signature, &user.Status}
}

func (user *User) CreateSqlRow() []any {
	return []any{user.Username, user.Email, user.HashedPwd.String, user.Salt.String, user.AvatarPath.String, user.BirthDate.Time, user.Bio.String, user.Signature.String}
}

// DbData interface implementation

func (user *User) Create() int {
	//TODO implement me
	panic("implement me")
	id := 1

	return id
}

func (user *User) Fetch(a any) {
	//TODO implement me
	panic("implement me")
}

func (user *User) GetId(a any) int {
	var login string
	switch a.(type) {
	case string:
		login = a.(string)
	default:
		// todo handle error
		return -1
	}
	id, err := db.GetUserByLogin(login)
	if err != nil {
		// todo handle error
		return -1
	}
	return id.Id
}

func (user *User) Exists(a any) bool {
	//TODO implement me
	panic("implement me")
}

func (user *User) Update(a any) {
	//TODO implement me
	panic("implement me")
}

// Category methods

// DbData interface implementation

func (c *Category) Create() int {
	//TODO implement me
	panic("implement me")
	id := 1

	return id
}

func (c *Category) Fetch(a any) {
	//TODO implement me
	panic("implement me")
}

func (c *Category) GetId(a any) int {
	//TODO implement me
	panic("implement me")
}

func (c *Category) Exists(a any) bool {
	//TODO implement me
	panic("implement me")
}

func (c *Category) Update(a any) {
	//TODO implement me
	panic("implement me")
}

// Tag methods

// DbData interface implementation

func (t *Tag) Create() int {
	//TODO implement me
	panic("implement me")
	id := 1

	return id
}

func (t *Tag) Fetch(a any) {
	//TODO implement me
	panic("implement me")
}

func (t *Tag) GetId(a any) int {
	//TODO implement me
	panic("implement me")
}

func (t *Tag) Exists(a any) bool {
	//TODO implement me
	panic("implement me")
}

func (t *Tag) Update(a any) {
	//TODO implement me
	panic("implement me")
}

// Thread methods

// DbData interface implementation

func (t *Thread) Create() int {
	//TODO implement me
	panic("implement me")
	id := 1

	return id
}

func (t *Thread) Fetch(a any) {
	//TODO implement me
	panic("implement me")
}

func (t *Thread) GetId(a any) int {
	//TODO implement me
	panic("implement me")
}

func (t *Thread) Exists(a any) bool {
	//TODO implement me
	panic("implement me")
}

func (t *Thread) Update(a any) {
	//TODO implement me
	panic("implement me")
}

// Post methods

// DbData interface implementation

func (p *Post) Create() int {
	//TODO implement me
	panic("implement me")
	id := 1

	return id
}

func (p *Post) Fetch(a any) {
	//TODO implement me
	panic("implement me")
}

func (p *Post) GetId(a any) int {
	//TODO implement me
	panic("implement me")
}

func (p *Post) Exists(a any) bool {
	//TODO implement me
	panic("implement me")
}

func (p *Post) Update(a any) {
	//TODO implement me
	panic("implement me")
}

// Friend methods

// DbData interface implementation

func (f *Friend) Create() int {
	//TODO implement me
	panic("implement me")
	id := 1

	return id
}

func (f *Friend) Fetch(a any) {
	//TODO implement me
	panic("implement me")
}

func (f *Friend) GetId(a any) int {
	//TODO implement me
	panic("implement me")
}

func (f *Friend) Exists(a any) bool {
	//TODO implement me
	panic("implement me")
}

func (f *Friend) Update(a any) {
	//TODO implement me
	panic("implement me")
}
