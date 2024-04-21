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

	// Making the query for a single row output
	row := db.QueryRow(getUserByIdQuery, id)

	// Scanning the result to put it in the user variable
	if err := row.Scan(user.GetSqlRows()...); err != nil {

		// Checking if the query returns nothing
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

	// Making the query for a single row output
	row := db.QueryRow(getUserByLoginQuery, login, login)

	// Scanning the result to put it in the user variable
	if err := row.Scan(user.GetSqlRows()...); err != nil {

		// Checking if the query returns nothing
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("getUserByLoginQuery %s: no such user", login)
		}
		return user, fmt.Errorf("getUserByLoginQuery %s: %v", login, err)
	}
	return user, nil
}

func UpdateUser(user models.User, updatedFields map[string]any) error {

	// Checking if there are info to update
	if len(updatedFields) == 0 {
		return fmt.Errorf("no field to update for user %s", user.Username)
	}

	// Setting the query's SET clauses variable
	var setClause string

	// Setting the slice variable holding the query's SET clauses values
	var args []interface{}

	// Looping through the updatedFields map to set the appropriate clauses
	// and values for the query
	for key, value := range updatedFields {
		setClause += key + " = ?, "
		args = append(args, value)
	}

	// Adding the Updated_at field to refresh it at the same time
	setClause += "Updated_at = CURRENT_TIMESTAMP"

	log.Printf("setClause: %s\n", setClause) // testing

	// Inserting the SET clauses in the query's body
	query := fmt.Sprintf(updateUserQuery, setClause)

	// Adding the user's Id to the values
	args = append(args, user.Id)

	log.Printf("query: %s\nargs: %v\n", query, args) // testing

	// Executing the query with the values
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("UpdateUser: %v", err)
	}

	log.Printf("Updated user with id: %d\n", user.Id) // testing

	return nil
}

func CreateUser(user models.User) error {

	// Executing the query to create the user
	_, err := db.Exec(createUserQuery, user.CreateSqlRow()...)
	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}

	log.Printf("Created user with username: %s\n", user.Username) // testing

	return nil
}

func DeleteUser(user models.User) error {

	// Retrieving the user's Id
	id := user.Id

	// Executing the query to delete the user
	_, err := db.Exec(deleteUserQuery, id)
	if err != nil {
		return fmt.Errorf("DeleteUser: %v", err)
	}

	log.Printf("Deleted user with id: %d\n", id) // testing

	return nil
}

func IsLogin(login string) (bool, error) {

	// Setting the variable holding the query's result
	var result int

	// Making the query for a single row output
	row := db.QueryRow(isLoginQuery, login, login)

	// Scanning the result to put it in the result variable
	if err := row.Scan(&result); err != nil {

		// Checking if the query returns nothing
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("IsLogin: %v", err)
	}

	log.Printf("Found %d users with %s as login\n", result, login) // testing

	return result > 0, nil
}
