package data

import (
	"ForumAPI/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) Insert(post Post) error {

	args := []any{post.Content, post.Author.ID}
	var parentPost, value string

	if post.IDParentPost != 0 {
		args = append(args, post.IDParentPost)
		parentPost = ", Id_parent_posts"
		value = ", ?"
	}

	query := fmt.Sprintf(`
		INSERT INTO posts (Content, Id_author%s)
		VALUES (?, ?%s);`, parentPost, value)

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
	postID, err := rs.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int(postID)
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	query = `
		SELECT Created_at, Version
		FROM posts
		WHERE Id_posts = ?;`

	err = tx.QueryRowContext(ctx, query, post.ID).Scan(&post.CreatedAt, &post.Version)
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

func (m PostModel) Get(id int) (*Post, error) {

	query := `
		SELECT m.Id_posts, m.Content, m.Created_at, m.Updated_at, m.Id_author, u.Username, m.Id_parent_posts, m.Id_threads, t.Title, m.Version
		FROM posts m
		INNER JOIN users u ON m.Id_author = u.Id_users
		INNER JOIN threads t ON m.Id_threads = t.Id_threads
		WHERE m.Id_posts = ?;`

	var post Post
	var parentPost sql.NullInt64

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Author.ID,
		&post.Author.Name,
		&parentPost,
		&post.Thread.ID,
		&post.Thread.Title,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if parentPost.Valid {
		post.IDParentPost = int(parentPost.Int64)
	}

	return &post, nil
}

func (m PostModel) GetPostsByAuthorID(id int) ([]Post, error) {

	query := `
		SELECT p.Id_posts, p.Content, p.Created_at, p.Updated_at, p.Id_threads, t.Title
		FROM posts p
		INNER JOIN forum.threads t on p.Id_threads = t.Id_threads
		WHERE p.Id_author = ?;`

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

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Thread.ID, &post.Thread.Title); err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return posts, nil
}

func (m PostModel) Update(post Post) error {

	query := `
		UPDATE posts 
		SET Content = ?, Id_author= ?, Id_parent_posts = ?, Version = Version + 1
		WHERE Id_posts = ? AND Version = ?;`

	args := []any{post.Content, post.Author.ID, post.IDParentPost, post.ID, post.Version}

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
		FROM posts
		WHERE Id_posts = ?;`

	err = tx.QueryRowContext(ctx, query, post.ID).Scan(&post.CreatedAt, &post.Version)
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

func (m PostModel) Delete(id int) error {

	query := `
		DELETE FROM posts
		WHERE Id_posts = ?;`

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

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	IDParentPost int `json:"id_parent_post"`
	Thread       struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"thread"`
	Version int `json:"version"`
}

func (post *Post) Validate(v *validator.Validator) {
	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 1_020, "content", "must not be more than 1020 bytes long")
	v.Check(post.Author.Name != "", "author.name", "must be provided")
	v.Check(len(post.Author.Name) <= 30, "author.name", "must not be more than 30 bytes long")
	v.Check(post.Thread.ID != 0, "post.thread.id", "must be provided")
	v.Check(post.Thread.Title != "", "post.thread.title", "must be provided")
	v.Check(len(post.Thread.Title) <= 125, "post.thread.title", "must not be more than 125 bytes long")
}
