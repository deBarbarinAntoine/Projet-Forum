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

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) Insert(category Category) error {

	args := []any{category.Name, category.Author.ID}
	var parentCategory, value string

	if category.ParentCategory.ID != 0 {
		args = append(args, category.ParentCategory.ID)
		parentCategory = ", Id_parent_categories"
		value = ", ?"
	}

	query := fmt.Sprintf(`
		INSERT INTO categories (Name, Id_author%s)
		VALUES (?, ?%s);`, parentCategory, value)

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
	categoryID, err := rs.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = int(categoryID)
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	query = `
		SELECT Created_at, Version
		FROM categories
		WHERE Id_categories = ?;`

	err = tx.QueryRowContext(ctx, query, category.ID).Scan(&category.CreatedAt, &category.Version)
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

func (m CategoryModel) Get(id int) (Category, error) {

	query := `
		SELECT Id_categories, Name, Id_parent_categories, Created_at, Version
		FROM categories
		WHERE Id_categories = ?;`

	var category Category

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.ParentCategory.ID,
		&category.CreatedAt,
		&category.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Category{}, ErrRecordNotFound
		default:
			return Category{}, err
		}
	}

	return category, nil
}

func (m CategoryModel) GetOwnedCategoriesByUserID(id int) ([]Category, error) {

	query := `
		SELECT Id_categories, Name, Created_at, Updated_at
		FROM categories
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

	var categoriesOwned []Category

	for rows.Next() {
		var categoryOwned Category
		if err := rows.Scan(&categoryOwned.ID, &categoryOwned.Name, &categoryOwned.CreatedAt, &categoryOwned.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		categoriesOwned = append(categoriesOwned, categoryOwned)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return categoriesOwned, nil
}

func (m CategoryModel) Update(category Category) error {

	query := `
		UPDATE categories 
		SET Name = ?, Id_author= ?, Id_parent_categories = ?, Version = Version + 1
		WHERE Id_categories = ? AND Version = ?;`

	args := []any{category.Name, category.Author.ID, category.ParentCategory.ID, category.ID, category.Version}

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
		FROM categories
		WHERE Id_categories = ?;`

	err = tx.QueryRowContext(ctx, query, category.ID).Scan(&category.CreatedAt, &category.Version)
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

func (m CategoryModel) Delete(id int) error {

	query := `
		DELETE FROM categories
		WHERE Id_categories = ?;`

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

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	ParentCategory struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"parent_category"`
	Version    int        `json:"version"`
	Categories []Category `json:"categories"`
	Threads    []Thread   `json:"threads"`
}

func (category *Category) Validate(v *validator.Validator) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 50, "name", "must not be more than 50 bytes long")
	v.Check(category.Author.Name != "", "author.name", "must be provided")
	v.Check(len(category.Author.Name) <= 30, "author.name", "must not be more than 30 bytes long")
	if category.ParentCategory.ID != 0 || category.ParentCategory.Name != "" {
		v.Check(category.ParentCategory.ID != 0, "parent_category.id", "must be provided")
		v.Check(category.ParentCategory.Name != "", "parent_category.name", "must be provided")
		v.Check(len(category.ParentCategory.Name) <= 50, "parent_category.name", "must not be more than 50 bytes long")
	}
}
