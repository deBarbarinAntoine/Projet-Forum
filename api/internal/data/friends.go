package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Friend struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (m UserModel) GetFriendsByUserID(id int) ([]Friend, error) {

	query := `
		SELECT users.Id_users, users.Username, f1.Status, f2.Status
		FROM users
		LEFT OUTER JOIN friends f1 ON ? = f1.Id_users_to
		LEFT OUTER JOIN friends f2 ON ? = f2.Id_users_from
		WHERE f1.Id_users_from = users.Id_users OR f2.Id_users_to = users.Id_users;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id, id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var friends []Friend

	for rows.Next() {
		var friend Friend
		var status1, status2 *string
		if err := rows.Scan(&friend.ID, &friend.Name, &status1, &status2); err != nil {
			log.Fatal(err)
		}
		switch {
		case status1 != nil:
			friend.Status = *status1
		case status2 != nil:
			friend.Status = *status2
		}
		friends = append(friends, friend)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return friends, nil
}
