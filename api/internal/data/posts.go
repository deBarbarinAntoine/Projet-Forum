package data

import (
	"ForumAPI/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Post struct {
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Author       User      `json:"author"`
	IDParentPost int       `json:"id_parent_post,omitempty"`
	Thread       struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"thread"`
	Reactions  map[string]int `json:"reactions,omitempty"`
	Popularity int            `json:"popularity,omitempty"`
	Version    int            `json:"version,omitempty"`
}

func (post *Post) Validate(v *validator.Validator) {
	v.StringCheck(post.Content, 2, 1_020, true, "content")
	v.StringCheck(post.Author.Name, 2, 70, true, "author.name")
	v.StringCheck(post.Thread.Title, 2, 125, true, "thread.title")
	v.Check(post.Thread.ID != 0, "post.thread.id", "must be provided")
}

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) Insert(post *Post) error {

	args := []any{post.Content, post.Author.ID, post.Thread.ID}
	var parentPost, value string

	if post.IDParentPost != 0 {
		args = append(args, post.IDParentPost)
		parentPost = ", Id_parent_posts"
		value = ", ?"
	}

	query := fmt.Sprintf(`
		INSERT INTO posts (Content, Id_author, Id_threads%s)
		VALUES (?, ?, ?%s);`, parentPost, value)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1452 {
				return ErrRecordNotFound
			}
		}
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

func (m PostModel) Get(search string, filters Filters) ([]*Post, Metadata, error) {

	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
	} else {
		search = "%"
	}

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), p.Id_posts, p.Content, p.Created_at, p.Updated_at, p.Id_author, u.Username, u.Avatar_path, p.Id_threads, t.Title
		FROM posts p
		INNER JOIN users u ON p.Id_author = u.Id_users
		INNER JOIN threads t ON p.Id_threads = t.Id_threads
		WHERE p.Content LIKE ?
		ORDER BY %s %s, Id_posts ASC
		LIMIT ? OFFSET ?;`, filters.sortColumn(), filters.sortDirection())

	args := []any{search, filters.limit, filters.offset}

	var posts []*Post

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, Metadata{}, ErrRecordNotFound
		default:
			return nil, Metadata{}, err
		}
	}

	var totalRecords int

	for rows.Next() {
		var post Post

		err := rows.Scan(
			&totalRecords,
			&post.ID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author.ID,
			&post.Author.Name,
			&post.Author.Avatar,
			&post.Thread.ID,
			&post.Thread.Title,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return posts, metadata, nil
}

func (m PostModel) GetByID(id int) (*Post, error) {

	query := `
		SELECT p.Id_posts, p.Content, p.Created_at, p.Updated_at, p.Id_author, u.Username, u.Avatar_path, p.Id_parent_posts, p.Id_threads, t.Title, p.Version
		FROM posts p
		INNER JOIN users u ON p.Id_author = u.Id_users
		INNER JOIN threads t ON p.Id_threads = t.Id_threads
		WHERE p.Id_posts = ?;`

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
		&post.Author.Avatar,
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

func (m PostModel) GetByAuthorID(id int) ([]Post, error) {

	query := `
		SELECT p.Id_posts, p.Content, p.Created_at, p.Updated_at, p.Id_threads, t.Title, p.Version
		FROM posts p
		INNER JOIN threads t on p.Id_threads = t.Id_threads
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
		if err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Thread.ID,
			&post.Thread.Title,
			&post.Version); err != nil {
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

func (m PostModel) GetByThread(id int) ([]*Post, error) {

	query := `
		SELECT p.Id_posts, p.Content, p.Created_at, p.Updated_at, p.Id_author, u.Username, u.Avatar_path, p.Version
		FROM posts p
		INNER JOIN users u on p.Id_author = u.Id_users
		WHERE p.Id_threads = ?;`

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

	var posts []*Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author.ID,
			&post.Author.Name,
			&post.Author.Avatar,
			&post.Version); err != nil {
			log.Fatal(err)
		}
		posts = append(posts, &post)
	}

	err = rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	if err = rows.Err(); err != nil {
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

func (m PostModel) GetReactions(posts []*Post) error {

	query := `
	SELECT p.Id_posts, e.Emoji, COUNT(*) AS Count
	FROM posts_users pu
	INNER JOIN posts p ON pu.Id_posts = p.Id_posts
	INNER JOIN (
		SELECT DISTINCT Emoji
		FROM posts_users
	) AS e ON pu.Emoji = e.Emoji
	WHERE p.Id_posts IN (?)
	GROUP BY p.Id_posts, e.Emoji
	ORDER BY p.Id_posts, e.Emoji;`

	var IDs []any

	for _, post := range posts {
		IDs = append(IDs, post.ID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	results := make(map[int]map[string]int)
	rows, err := stmt.QueryContext(ctx, IDs...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var postID int
		var emoji string
		var count int
		err := rows.Scan(&postID, &emoji, &count)
		if err != nil {
			return err
		}

		if _, ok := results[postID]; !ok {
			results[postID] = make(map[string]int)
		}

		results[postID][emoji] = count
	}

	for _, post := range posts {
		post.Reactions = results[post.ID]
		for _, i := range results[post.ID] {
			post.Popularity += i
		}
	}

	return nil
}

func (m PostModel) React(user *User, id int, reaction string) error {

	query := `
		INSERT INTO posts_users (Id_users, Id_posts, Emoji)
		VALUES (?, ?, ?);`

	args := []any{user.ID, id, reaction}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				return ErrDuplicateEntry
			}
			if mySQLError.Number == 1452 {
				return ErrRecordNotFound
			}
		default:
			return err
		}
	}

	return nil
}

func (m PostModel) UpdateReaction(user *User, id int, reaction string) error {

	query := `
		UPDATE posts_users
		SET Emoji = ?
		WHERE Id_users = ? AND Id_posts = ?;`

	args := []any{reaction, user.ID, id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1452 {
				return ErrRecordNotFound
			}
		default:
			return err
		}
	}

	return nil
}

func (m PostModel) RemoveReaction(user *User, id int) error {

	query := `
		DELETE FROM posts_users
		WHERE Id_users = ? AND Id_posts = ?;`

	args := []any{user.ID, id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1452 {
				return ErrRecordNotFound
			}
		default:
			return err
		}
	}

	return nil
}
