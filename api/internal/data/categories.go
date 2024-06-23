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

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) Insert(category *Category) error {

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

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
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
	categoryID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = int(categoryID)
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	query = `
		SELECT c.Created_at, pc.Name, c.Version
		FROM categories c
		LEFT JOIN categories pc ON c.Id_parent_categories = pc.Id_categories
		WHERE c.Id_categories = ?;`

	var parentCategoryName sql.NullString

	err = tx.QueryRowContext(ctx, query, category.ID).Scan(&category.CreatedAt, &parentCategoryName, &category.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	if parentCategoryName.Valid {
		category.ParentCategory.Name = parentCategoryName.String
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m CategoryModel) Get(search string, filters Filters) ([]*Category, Metadata, error) {

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), c.Id_categories, c.Name, c.Id_author, u.Username, c.Id_parent_categories, pc.Name, c.Created_at, c.Updated_at, c.Version
		FROM categories c
		INNER JOIN users u ON u.Id_users = c.Id_author
		LEFT OUTER JOIN categories pc ON pc.Id_categories = c.Id_parent_categories
		WHERE c.Name LIKE ?
		ORDER BY %s %s, Id_Categories ASC
		LIMIT ? OFFSET ?;`, filters.sortColumn(), filters.sortDirection())

	args := []any{search, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
	} else {
		search = "%"
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	var totalRecords int
	var categories []*Category

	for rows.Next() {

		var parentID sql.NullInt64
		var parentName sql.NullString
		var category Category

		err = rows.Scan(
			&totalRecords,
			&category.ID,
			&category.Name,
			&category.Author.ID,
			&category.Author.Name,
			&parentID,
			&parentName,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		if parentID.Valid {
			category.ParentCategory.ID = int(parentID.Int64)
		}
		if parentName.Valid {
			category.ParentCategory.Name = parentName.String
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return categories, metadata, nil
}

func (m CategoryModel) GetByID(id int) (Category, error) {

	query := `
		SELECT c.Id_categories, c.Name, c.Id_author, u.Username, c.Id_parent_categories, pc.Name, c.Created_at, c.Updated_at, c.Version
		FROM categories c
		INNER JOIN users u ON u.Id_users = c.Id_author
		INNER JOIN categories pc ON pc.Id_categories = c.Id_parent_categories
		WHERE c.Id_categories = ?;`

	var category Category

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var parentID sql.NullInt64
	var parentName sql.NullString

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Author.ID,
		&category.Author.Name,
		&parentID,
		&parentName,
		&category.CreatedAt,
		&category.UpdatedAt,
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

	if parentID.Valid {
		category.ParentCategory.ID = int(parentID.Int64)
	}
	if parentName.Valid {
		category.ParentCategory.Name = parentName.String
	}

	return category, nil
}

func (m CategoryModel) GetByParentID(id int) ([]Category, error) {

	query := `
		SELECT c.Id_categories, c.Name, c.Id_author, u.Username, c.Created_at, c.Updated_at
		FROM categories c
		INNER JOIN users u ON u.Id_users = c.Id_author
		WHERE c.Id_parent_categories = ?;`

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

	var categories []Category

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Author.ID, &category.Author.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return categories, nil
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

func (m CategoryModel) Update(category *Category) error {

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
		var mySQLError *mysql.MySQLError
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
		return ErrEditConflict
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
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"parent_category,omitempty"`
	Version    int        `json:"version,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Threads    []Thread   `json:"threads,omitempty"`
}

func (category *Category) Validate(v *validator.Validator) {
	v.StringCheck(category.Name, 2, 70, true, "name")
	v.StringCheck(category.Author.Name, 2, 70, true, "author.name")

	if category.ParentCategory.ID != 0 || category.ParentCategory.Name != "" {
		v.Check(category.ParentCategory.ID != 0, "parent_category.id", "must be provided")
		v.StringCheck(category.ParentCategory.Name, 2, 70, true, "parent_category.name")
	}
}
