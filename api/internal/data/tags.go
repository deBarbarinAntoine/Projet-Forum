package data

import (
	"ForumAPI/internal/validator"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type TagModel struct {
	DB *sql.DB
}

func (m TagModel) Insert(tag Tag) error {

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

	err = tx.QueryRowContext(ctx, query, tag.ID).Scan(&tag.CreatedAt, &tag.Version)
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

func (m TagModel) Get(id int) (Tag, error) {

	query := `
		SELECT Id_tags, Name, Created_at, Id_author, Version
		FROM tags
		WHERE Id_tags = ?;`

	var tag Tag

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.CreatedAt,
		&tag.Author.ID,
		&tag.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Tag{}, ErrRecordNotFound
		default:
			return Tag{}, err
		}
	}

	return tag, nil
}

func (m TagModel) GetOwnedTagsByUserID(id int) ([]Tag, error) {

	query := `
		SELECT Id_tags, Name, Created_at, Updated_at
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
		if err := rows.Scan(&tagOwned.ID, &tagOwned.Name, &tagOwned.CreatedAt, &tagOwned.UpdatedAt); err != nil {
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

func (m TagModel) GetFollowingTagsByUserID(id int) ([]struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}, error) {

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

	var followingTags []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	for rows.Next() {
		var followingTag struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
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

func (m TagModel) Update(tag Tag) error {

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

	if err = tx.Commit(); err != nil {
		return err
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

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Version int      `json:"version"`
	Threads []Thread `json:"threads"`
}

func (tag *Tag) Validate(v *validator.Validator) {
	v.Check(tag.Name != "", "name", "must be provided")
	v.Check(len(tag.Name) <= 50, "name", "must not be more than 50 bytes long")
	v.Check(tag.Author.Name != "", "author.name", "must be provided")
	v.Check(len(tag.Author.Name) <= 30, "author.name", "must not be more than 30 bytes long")
}
