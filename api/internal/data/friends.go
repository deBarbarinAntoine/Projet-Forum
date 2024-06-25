package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

type Friend struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type friendStatus struct {
	Pending  string
	Accepted string
	Rejected string
}

var FriendStatus = friendStatus{
	Pending:  "pending",
	Accepted: "accepted",
	Rejected: "rejected",
}

func (m UserModel) RequestFriend(user *User, idTo int) error {

	query := `
		INSERT INTO friends (
		    Id_users_from,
		    Id_users_to,
		    Status)
		VALUES (
		    Id_users_from = ?,
		    Id_users_to = ?,
		    Status = ?);`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var mySQLError *mysql.MySQLError

	_, err := m.DB.ExecContext(ctx, query, user.ID, idTo, FriendStatus.Pending)
	if err != nil {
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Id_users") {
					return ErrDuplicateFriend
				}
			}
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) AcceptFriend(idFrom int, user *User) error {

	query := `
		UPDATE friends SET Status = ?
		WHERE Id_users_from = ? AND Id_users_to = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, FriendStatus.Accepted, idFrom, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) RejectFriend(idFrom int, user *User) error {

	query := `
		UPDATE friends SET Status = ?
		WHERE Id_users_from = ? AND Id_users_to = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, FriendStatus.Rejected, idFrom, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) RemoveFriend(user *User, id int) error {

	query := `
		DELETE f1, f2 FROM friends f1, friends f2
		WHERE (f1.Id_users_from = ? AND f1.Id_users_to = ?) OR (f2.Id_users_to = ? AND f2.Id_users_from = ?);`

	args := []any{user.ID, id, user.ID, id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, query, args)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetFriendsByUserID(id int) (sentFriends, receivedFriends []Friend, err error) {

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
			return nil, nil, ErrRecordNotFound
		default:
			return nil, nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var friend Friend
		var status1, status2 *string
		if err := rows.Scan(&friend.ID, &friend.Name, &status1, &status2); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, nil, ErrRecordNotFound
			default:
				return nil, nil, err
			}
		}
		switch {
		case status1 != nil:
			friend.Status = *status1
			sentFriends = append(sentFriends, friend)
		case status2 != nil:
			friend.Status = *status2
			receivedFriends = append(receivedFriends, friend)
		}
	}
	rerr := rows.Close()
	if rerr != nil {
		return nil, nil, rerr
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return sentFriends, receivedFriends, nil
}
