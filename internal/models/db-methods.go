package models

func (user *User) GetSqlRows() []any {
	return []any{&user.Id, &user.Username, &user.Email, &user.HashedPwd, &user.Salt, &user.AvatarPath, &user.Role, &user.BirthDate, &user.CreatedAt, &user.UpdatedAt, &user.VisitedAt, &user.Bio, &user.Signature, &user.Status}
}

func (user *User) CreateSqlRow() []any {
	return []any{user.Username, user.Email, user.HashedPwd.String, user.Salt.String, user.AvatarPath.String, user.BirthDate.Time, user.Bio.String, user.Signature.String}
}

//func (user *User) UpdateTo(updatedFields map[string]any) {
//
//}
