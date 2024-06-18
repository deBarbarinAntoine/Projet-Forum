package data

import (
	"ForumAPI/internal/validator"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type ThreadModel struct {
	DB *sql.DB
}

func (m ThreadModel) Insert(thread Thread) error {

	args := []any{thread.Title, thread.Description, thread.IsPublic, thread.Status, thread.Author.ID, thread.Category.ID}

	query := `
		INSERT INTO threads (Title, Description, Is_public, Status, Id_author, Id_categories)
		VALUES (?, ?, ?, ?, ?, ?);`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rs, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	threadID, err := rs.LastInsertId()
	if err != nil {
		return err
	}
	thread.ID = int(threadID)
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	query = `
		SELECT Created_at, Version
		FROM threads
		WHERE Id_threads = ?;`

	err = tx.QueryRowContext(ctx, query, thread.ID).Scan(&thread.CreatedAt, &thread.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m ThreadModel) Get(id int) (Thread, error) {

	query := `
		SELECT Id_threads, Title, Description, Is_public, Created_at, Status, Id_author, Id_categories, Version
		FROM threads
		WHERE Id_threads = ?;`

	var thread Thread

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Description,
		&thread.IsPublic,
		&thread.CreatedAt,
		&thread.Status,
		&thread.Author.ID,
		&thread.Category.ID,
		&thread.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Thread{}, ErrRecordNotFound
		default:
			return Thread{}, err
		}
	}

	return thread, nil
}

func (m ThreadModel) GetOwnedThreadsByUserID(id int) ([]Thread, error) {

	query := `
		SELECT Id_threads, Title, Description, Is_public, Created_at, Updated_at, Status
		FROM threads
		WHERE Id_author = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var threadsOwned []Thread

	for rows.Next() {
		var threadOwned Thread
		if err := rows.Scan(&threadOwned.ID, &threadOwned.Title, &threadOwned.Description, &threadOwned.IsPublic, &threadOwned.CreatedAt, &threadOwned.UpdatedAt, &threadOwned.Status); err != nil {
			log.Fatal(err)
		}
		threadsOwned = append(threadsOwned, threadOwned)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return threadsOwned, nil
}

func (m ThreadModel) GetFavoriteThreadsByUserID(id int) ([]struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}, error) {

	query := `
		SELECT tu.Id_threads, t.Title
		FROM threads_users tu
		INNER JOIN threads t ON tu.Id_threads = t.Id_threads
		WHERE Id_users = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var favoriteThreads []struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}

	for rows.Next() {
		var favoriteThread struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		}
		if err := rows.Scan(&favoriteThread.ID, &favoriteThread.Title); err != nil {
			log.Fatal(err)
		}
		favoriteThreads = append(favoriteThreads, favoriteThread)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return favoriteThreads, nil
}

func (m ThreadModel) Update(thread Thread) error {

	query := `
		UPDATE threads 
		SET Title = ?, Description = ?, Is_public = ?, Status = ?, Id_author = ?, Id_categories = ?, Version = Version + 1
		WHERE Id_threads = ? AND Version = ?;`

	args := []any{thread.Title, thread.Description, thread.IsPublic, thread.Status, thread.Author.ID, thread.Category.ID, thread.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rs, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	query = `
		SELECT Created_at, Version
		FROM threads
		WHERE Id_threads = ?;`

	err = tx.QueryRowContext(ctx, query, thread.ID).Scan(&thread.CreatedAt, &thread.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m ThreadModel) Delete(id int) error {

	query := `
		DELETE FROM threads
		WHERE Id_threads = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type Thread struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	Author      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Category struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Version int    `json:"version"`
	Posts   []Post `json:"posts"`
	Tags    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

func (thread *Thread) Validate(v *validator.Validator) {
	v.Check(thread.Title != "", "title", "must be provided")
	v.Check(len(thread.Title) <= 125, "title", "must not be more than 125 bytes long")
	v.Check(thread.Description != "", "description", "must be provided")
	v.Check(len(thread.Description) <= 1_020, "description", "must not be more than 1020 bytes long")
	v.Check(thread.IsPublic, "is_public", "must be provided")
	v.Check(thread.Author.Name != "", "author.name", "must be provided")
	v.Check(len(thread.Author.Name) <= 30, "author.name", "must not be more than 30 bytes long")
	v.Check(thread.Category.ID != 0, "parent_category.id", "must be provided")
	v.Check(thread.Category.Name != "", "parent_category.name", "must be provided")
	v.Check(len(thread.Category.Name) <= 50, "parent_category.name", "must not be more than 50 bytes long")
}