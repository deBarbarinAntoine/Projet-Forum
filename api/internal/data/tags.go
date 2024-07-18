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

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Author    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Version    int      `json:"version,omitempty"`
	Popularity int      `json:"popularity,omitempty"`
	Threads    []Thread `json:"threads,omitempty"`
}

func (tag *Tag) Validate(v *validator.Validator) {
	v.StringCheck(tag.Name, 2, 50, true, "name")
	v.StringCheck(tag.Author.Name, 2, 70, true, "author.name")
}

type TagModel struct {
	DB *sql.DB
}

func (m TagModel) Insert(tag *Tag) error {

	args := []any{tag.Name, tag.Author.ID}

	query := `
		INSERT INTO tags (Name, Id_author)
		VALUES (?, ?);`

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
	tagID, err := rs.LastInsertId()
	if err != nil {
		return err
	}
	tag.ID = int(tagID)
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	query = `
		SELECT Created_at, Version
		FROM tags
		WHERE Id_tags = ?;`

	var mySQLError *mysql.MySQLError

	err = tx.QueryRowContext(ctx, query, tag.ID).Scan(&tag.CreatedAt, &tag.Version)
	if err != nil {
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Name") {
					return ErrDuplicateName
				}
			}
		default:
			return err
		}
	}

	if len(tag.Threads) > 0 {

		var ids []any
		for _, thread := range tag.Threads {
			ids = append(ids, thread.ID)
		}

		value := fmt.Sprintf(`(?, %d),`, tag.ID)
		values := strings.Repeat(value, len(tag.Threads))
		values = values[:len(values)-1]

		query = fmt.Sprintf(`
		INSERT INTO threads_tags (Id_threads, Id_tags)
		VALUES %s;`, values)

		rs, err = tx.ExecContext(ctx, query, ids...)
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
		tagID, err := rs.LastInsertId()
		if err != nil {
			return err
		}
		tag.ID = int(tagID)
		rowsAffected, err := rs.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return ErrRecordNotFound
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m TagModel) GetByID(id int) (*Tag, error) {

	query := `
		SELECT t.Id_tags, t.Name, t.Created_at, t.Updated_at, t.Id_author, u.Username, t.Version
		FROM tags t
		INNER JOIN users u on t.Id_author = u.Id_users
		WHERE t.Id_tags = ?;`

	var tag Tag

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&tag.Author.ID,
		&tag.Author.Name,
		&tag.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tag, nil
}

func (m TagModel) GetByAuthorID(id int) ([]Tag, error) {

	query := `
		SELECT Id_tags, Name, Created_at, Updated_at, Version
		FROM tags
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

	var tagsOwned []Tag

	for rows.Next() {
		var tagOwned Tag
		if err := rows.Scan(&tagOwned.ID, &tagOwned.Name, &tagOwned.CreatedAt, &tagOwned.UpdatedAt, &tagOwned.Version); err != nil {
			log.Fatal(err)
		}
		tagsOwned = append(tagsOwned, tagOwned)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tagsOwned, nil
}

func (m TagModel) GetPopularity(id int) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM tags_users
		WHERE Id_tags = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var popularity int

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&popularity)
	if err != nil {
		return 0, err
	}

	return popularity, nil
}

func (m TagModel) GetPopular() ([]*Tag, error) {

	query := `
		SELECT t.Id_tags, t.Name, t.Id_author, u.Username, t.Created_at, t.Updated_at, (SELECT COUNT(*)
																						FROM tags_users tu
																						WHERE tu.Id_tags = t.Id_tags) AS popularity
		FROM tags t
		INNER JOIN users u on t.Id_author = u.Id_users
		ORDER BY popularity DESC, Id_tags ASC
		LIMIT 10;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag

	for rows.Next() {
		var tag Tag

		err := rows.Scan(&tag.ID, &tag.Name, &tag.Author.ID, &tag.Author.Name, &tag.CreatedAt, &tag.UpdatedAt, &tag.Popularity)
		if err != nil {
			return nil, err
		}

		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m TagModel) Get(search string, filters Filters) ([]*Tag, Metadata, error) {

	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
	} else {
		search = "%"
	}

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), t.Id_tags, t.Name, t.Created_at, t.Updated_at, t.Id_author, u.Username, t.Version
		FROM tags t
		INNER JOIN users u ON t.Id_author = u.Id_users
		WHERE t.Name LIKE ?
		ORDER BY %s %s, Id_tags ASC
		LIMIT ? OFFSET ?;`, filters.sortColumn(), filters.sortDirection())

	args := []any{search, filters.limit(), filters.offset()}

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
	defer rows.Close()

	var tags []*Tag
	var totalRecords int

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(
			&totalRecords,
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&tag.Author.ID,
			&tag.Author.Name,
			&tag.Version); err != nil {
			log.Fatal(err)
		}
		tags = append(tags, &tag)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tags, metadata, nil
}

func (m TagModel) GetByThread(id int) ([]Tag, error) {

	query := `
		SELECT tt.Id_tags, t.Name, t.Id_author, u.Username
		FROM threads_tags tt
		INNER JOIN tags t ON tt.Id_tags = t.Id_tags
		INNER JOIN users u ON t.Id_author = u.Id_users
		WHERE tt.Id_threads = ?;`

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

	var tags []Tag

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Author.ID, &tag.Author.Name); err != nil {
			log.Fatal(err)
		}
		tags = append(tags, tag)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tags, nil
}

func (m TagModel) GetByFollowingUserID(id int) ([]Tag, error) {

	query := `
		SELECT tu.Id_tags, t.Name
		FROM tags_users tu
		INNER JOIN tags t ON tu.Id_tags = t.Id_tags
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

	var followingTags []Tag

	for rows.Next() {
		var followingTag Tag
		if err := rows.Scan(&followingTag.ID, &followingTag.Name); err != nil {
			log.Fatal(err)
		}
		followingTags = append(followingTags, followingTag)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return followingTags, nil
}

func (m TagModel) Update(tag *Tag, removeThreads []int) error {

	query := `
		UPDATE tags 
		SET Name = ?, Id_author = ?, Version = Version + 1
		WHERE Id_tags = ? AND Version = ?;`

	args := []any{tag.Name, tag.Author.ID, tag.ID, tag.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var mySQLError *mysql.MySQLError

	rs, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Name") {
					return ErrDuplicateName
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
		SELECT Created_at, Version
		FROM tags
		WHERE Id_tags = ?;`

	err = tx.QueryRowContext(ctx, query, tag.ID).Scan(&tag.CreatedAt, &tag.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	if len(tag.Threads) > 0 {

		var ids []any
		for _, thread := range tag.Threads {
			ids = append(ids, thread.ID)
		}

		value := fmt.Sprintf(`(?, %d),`, tag.ID)
		values := strings.Repeat(value, len(tag.Threads))
		values = values[:len(values)-1]

		query = fmt.Sprintf(`
		INSERT INTO threads_tags (Id_threads, Id_tags)
		VALUES %s;`, values)

		rs, err = tx.ExecContext(ctx, query, ids...)
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
		rowsAffected, err := rs.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return ErrRecordNotFound
		}
	}

	if len(removeThreads) > 0 {

		ids := []any{tag.ID, removeThreads}

		query = `
		DELETE FROM threads_tags
		WHERE Id_tags = ? AND Id_threads IN (?);`

		rs, err = tx.ExecContext(ctx, query, ids...)
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
		rowsAffected, err := rs.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return ErrRecordNotFound
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m TagModel) Follow(user *User, id int) error {

	query := `
		INSERT INTO tags_users (Id_users, Id_tags)
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
			if mySQLError.Number == 1452 {
				return ErrRecordNotFound
			}
		default:
			return err
		}
	}

	var followingTag Tag
	followingTag.ID = id

	query = `
		SELECT Name
		FROM tags
		WHERE Id_tags = ?;`

	err = tx.QueryRowContext(ctx, query, id).Scan(&followingTag.Name)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	if user.FollowingTags == nil {
		user.FollowingTags = make([]Tag, 0)
	}

	user.FollowingTags = append(user.FollowingTags, followingTag)

	return nil
}

func (m TagModel) Unfollow(user *User, id int) error {

	query := `
		DELETE FROM tags_users 
		WHERE Id_users = ? AND Id_tags = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{user.ID, id}

	res, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return ErrRecordNotFound
	}

	if user.FollowingTags != nil {
		for i, followingTag := range user.FollowingTags {
			if followingTag.ID == id {
				user.FollowingTags = append(user.FollowingTags[:i], user.FollowingTags[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (m TagModel) Delete(id int) error {

	query := `
		DELETE FROM tags
		WHERE Id_tags = ?;`

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
