package db

import (
	"Projet-Forum/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func GetAllUsers() ([]models.User, error) {

	// A users slice to hold data from returned rows.
	var users []models.User

	rows, err := db.Query(getAllUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var user models.User
		if err := rows.Scan(user.GetSqlRows()...); err != nil {
			return nil, fmt.Errorf("GetAllUsers: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	return users, nil
}

func GetUserById(id int) (models.User, error) {

	// A user to hold data from the returned row.
	var user models.User

	row := db.QueryRow(getUserByIdQuery, id)
	if err := row.Scan(user.GetSqlRows()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("GetUserById %d: no such user", id)
		}
		return user, fmt.Errorf("GetUserById %d: %v", id, err)
	}
	return user, nil
}

func GetUserByLogin(login string) (models.User, error) {

	// A user to hold data from the returned row.
	var user models.User

	row := db.QueryRow(getUserByLoginQuery, login, login)
	if err := row.Scan(user.GetSqlRows()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("getUserByLoginQuery %s: no such user", login)
		}
		return user, fmt.Errorf("getUserByLoginQuery %s: %v", login, err)
	}
	return user, nil
}

func UpdateUser(user models.User, updatedFields map[string]any) error {
	var setClause string
	var args []interface{}
	for key, value := range updatedFields {
		setClause += key + " = ?, "
		args = append(args, value)
	}
	setClause += "Updated_at = CURRENT_TIMESTAMP"
	log.Printf("setClause: %s\n", setClause)
	query := fmt.Sprintf(updateUserQuery, setClause)
	args = append(args, user.Id)
	log.Printf("query: %s\nargs: %v\n", query, args)
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("UpdateUser: %v", err)
	}
	log.Printf("Updated user with id: %d\n", user.Id)
	return nil
}

func CreateUser(user models.User) error {
	_, err := db.Exec(createUserQuery, user.CreateSqlRow()...)
	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}
	log.Printf("Created user with username: %s\n", user.Username)
	return nil
}

func DeleteUser(user models.User) error {
	id := user.Id
	_, err := db.Exec(deleteUserQuery, id)
	if err != nil {
		return fmt.Errorf("DeleteUser: %v", err)
	}
	log.Printf("Deleted user with id: %d\n", id)
	return nil
}

func IsLogin(login string) (bool, error) {
	var result int
	row := db.QueryRow(isLoginQuery, login, login)
	if err := row.Scan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("IsLogin: %v", err)
	}
	log.Printf("Found %d users with %s as login\n", result, login)
	return result > 0, nil
}
