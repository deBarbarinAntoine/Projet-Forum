package data

import (
	"ForumAPI/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

type ThreadModel struct {
	DB *sql.DB
}

type threadStatus struct {
	Active   string
	Archived string
	Hidden   string
}

var (
	ThreadStatus = threadStatus{
		Active:   "active",
		Archived: "archived",
		Hidden:   "hidden",
	}
	permittedStatuses = []string{ThreadStatus.Active, ThreadStatus.Archived, ThreadStatus.Hidden}
)

func (m ThreadModel) Insert(thread *Thread) error {

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
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Title") {
					return ErrDuplicateTitle
				}
			}
		default:
			return err
		}
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
		SELECT t.Created_at, c.Name, t.Version
		FROM threads t
		INNER JOIN categories c ON t.Id_categories = c.Id_categories
		WHERE t.Id_threads = ?;`

	err = tx.QueryRowContext(ctx, query, thread.ID).Scan(&thread.CreatedAt, &thread.Category.Name, &thread.Version)
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

func (m ThreadModel) Get(search string, filters Filters) ([]*Thread, Metadata, error) {

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), t.Id_threads, t.Title, t.Description, t.Is_public, t.Created_at, t.Updated_at, t.Id_author, u.Username, t.Id_categories, c.Name, t.Status
		FROM threads t
		INNER JOIN users u ON t.Id_author = u.Id_users
		INNER JOIN categories c ON t.Id_categories = c.Id_categories
		WHERE t.Title LIKE ? OR t.Description LIKE ?
		ORDER BY %s %s, Id_threads ASC
		LIMIT ? OFFSET ?;`, filters.sortColumn(), filters.sortDirection())

	args := []any{search, search, filters.limit, filters.offset}

	var threads []*Thread

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
	} else {
		search = "%"
	}

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
		var thread Thread

		err := rows.Scan(
			&totalRecords,
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.IsPublic,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.Author.ID,
			&thread.Author.Name,
			&thread.Category.ID,
			&thread.Category.Name,
			&thread.Status,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		threads = append(threads, &thread)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return threads, metadata, nil
}

func (m ThreadModel) GetByID(id int) (*Thread, error) {

	query := `
		SELECT Id_threads, Title, Description, Is_public, Created_at, Updated_at, Status, Id_author, Id_categories, Version
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
		&thread.UpdatedAt,
		&thread.Status,
		&thread.Author.ID,
		&thread.Category.ID,
		&thread.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &thread, nil
}

func (m ThreadModel) GetByCategory(id int) ([]Thread, error) {

	query := `
		SELECT t.Id_threads, t.Title, t.Description, t.Is_public, t.Created_at, t.Updated_at, t.Id_author, u.Username, t.Status
		FROM threads t
		INNER JOIN forum.users u on t.Id_author = u.Id_users
		WHERE t.Id_categories = ?;`

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

	var threads []Thread

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Description,
			&thread.IsPublic,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.Author.ID,
			&thread.Author.Name,
			&thread.Status); err != nil {
			log.Fatal(err)
		}
		threads = append(threads, thread)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return threads, nil
}

func (m ThreadModel) GetOwnedThreadsByUserID(id int) ([]Thread, error) {

	query := `
		SELECT Id_threads, Title, Description, Is_public, Created_at, Updated_at, Status, Id_categories, Version
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
		if err := rows.Scan(&threadOwned.ID, &threadOwned.Title, &threadOwned.Description, &threadOwned.IsPublic, &threadOwned.CreatedAt, &threadOwned.UpdatedAt, &threadOwned.Status, &threadOwned.Category.ID, &threadOwned.Version); err != nil {
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
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Title") {
					return ErrDuplicateTitle
				}
			}
		default:
			return err
		}
	}
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	query = `
		SELECT t.Created_at, c.Name, t.Version
		FROM threads t
		INNER JOIN categories c ON t.Id_categories = c.Id_categories
		WHERE t.Id_threads = ?;`

	err = tx.QueryRowContext(ctx, query, thread.ID).Scan(&thread.CreatedAt, &thread.Category.Name, &thread.Version)
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

func (m ThreadModel) GetPopularity(id int) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM threads_users
		WHERE Id_threads = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var popularity int

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&popularity)
	if err != nil {
		return 0, err
	}

	return popularity, nil
}

func (m ThreadModel) AddToFavorites(user *User, id int) error {

	query := `
		INSERT INTO threads_users (Id_users, Id_threads)
		VALUES (?, ?);`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	args := []any{user.ID, id}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				return ErrDuplicateEntry
			}
		default:
			return err
		}
	}

	var favoriteThread struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	favoriteThread.ID = id

	query = `
		SELECT Title
		FROM threads
		WHERE Id_threads = ?;`

	err = tx.QueryRowContext(ctx, query, id).Scan(&favoriteThread.Title)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	if user.FavoriteThreads == nil {
		user.FavoriteThreads = make([]struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		}, 0)
	}

	user.FavoriteThreads = append(user.FavoriteThreads, favoriteThread)

	return nil
}

func (m ThreadModel) RemoveFromFavorites(user *User, id int) error {

	query := `
		DELETE FROM threads_users 
		WHERE Id_users = ? AND Id_threads = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{user.ID, id}

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if user.FavoriteThreads != nil {
		for i, favoriteThread := range user.FavoriteThreads {
			if favoriteThread.ID == id {
				user.FavoriteThreads = append(user.FavoriteThreads[:i], user.FavoriteThreads[i+1:]...)
				break
			}
		}
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
	Version    int    `json:"version,omitempty"`
	Popularity int    `json:"popularity"`
	Posts      []Post `json:"posts,omitempty"`
	Tags       []Tag  `json:"tags,omitempty"`
}

func (thread *Thread) Validate(v *validator.Validator) {
	v.StringCheck(thread.Title, 2, 125, true, "title")
	v.StringCheck(thread.Description, 2, 1_020, true, "description")
	v.Check(validator.PermittedValue(thread.Status, permittedStatuses...), "status", "must be a permitted value")
	v.StringCheck(thread.Author.Name, 2, 70, true, "author.name")
	v.StringCheck(thread.Category.Name, 2, 70, true, "parent_category.name")
	v.Check(thread.Category.ID != 0, "parent_category.id", "must be provided")
}
