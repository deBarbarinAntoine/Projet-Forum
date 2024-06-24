package data

import (
	"ForumAPI/internal/validator"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
	MinAge            = 13
)

var (
	AnonymousUser = &User{}
	UserStatus    = &userStatus{
		Activated:  "activated",
		ToConfirm:  "to_confirm",
		Blocked:    "blocked",
		Client:     "client",
		HostSecret: "host_secret",
	}
	UserRole = &userRole{
		Secret:    "host_secret",
		Client:    "client",
		Admin:     "admin",
		Moderator: "moderator",
		Normal:    "normal",
	}
)

type userStatus struct {
	Activated  string
	ToConfirm  string
	Blocked    string
	Client     string
	HostSecret string
}

type userRole struct {
	Secret    string
	Client    string
	Admin     string
	Moderator string
	Normal    string
}

type User struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      password  `json:"-"`
	Role          string    `json:"role"`
	BirthDate     time.Time `json:"birth_date"`
	Bio           string    `json:"bio,omitempty"`
	Signature     string    `json:"signature,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Status        string    `json:"status"`
	Version       int       `json:"-"`
	FollowingTags []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"following_tags,omitempty"`
	FavoriteThreads []struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"favorite_threads,omitempty"`
	CategoriesOwned []Category `json:"categories_owned,omitempty"`
	TagsOwned       []Tag      `json:"tags_owned,omitempty"`
	ThreadsOwned    []Thread   `json:"threads_owned,omitempty"`
	Posts           []Post     `json:"posts,omitempty"`
	Friends         []Friend   `json:"friends,omitempty"`
	Invitations     struct {
		Received []Friend `json:"received"`
		Sent     []Friend `json:"sent"`
	} `json:"invitations,omitempty"`
}

func (u *User) IsActivated() bool {
	return u.Status == UserStatus.Activated
}

func (u *User) IsToConfirm() bool {
	return u.Status == UserStatus.ToConfirm
}

func (u *User) IsBlocked() bool {
	return u.Status == UserStatus.Blocked
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) NoLogin() {
	u.Password = password{
		plaintext: nil,
		hash:      []byte("no-login"),
	}
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plainTextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.StringCheck(password, MinPasswordLength, MaxPasswordLength, true, "password")
}

func ValidateNewPassword(v *validator.Validator, newPassword, confirmationPassword string) {
	v.StringCheck(newPassword, MinPasswordLength, MaxPasswordLength, true, "new_password")
	v.Check(confirmationPassword != "", "confirmation_password", "must be provided")
	v.Check(newPassword == confirmationPassword, "confirmation_password", "must be the same")
}

func (u *User) Validate(v *validator.Validator) {
	v.StringCheck(u.Name, 2, 70, true, "name")
	v.Check(u.BirthDate.AddDate(MinAge, 0, 0).Before(time.Now()), "birth", fmt.Sprintf("must be at least %d years old", MinAge))

	ValidateEmail(v, u.Email)

	if u.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *u.Password.plaintext)
	}

	if u.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {

	args := []any{user.Name, user.Email, user.Password.hash, user.Status}

	var role, roleValue string
	if user.Role != "" {
		role = ", Role"
		roleValue = ", ?"
		args = append(args, user.Role)
	}

	query := fmt.Sprintf(`
		INSERT INTO users (Username, Email, Hashed_password, Status%s)
		VALUES (?, ?, ?, ?%s);`, role, roleValue)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Email") {
					return ErrDuplicateEmail
				}
				if strings.Contains(mySQLError.Message, "Username") {
					return ErrDuplicateUsername
				}
			}
		default:
			return err
		}
	}

	query = `
		SELECT Id_users, Created_at, Version
		FROM users
		WHERE Username = ?;`

	err = tx.QueryRowContext(ctx, query, user.Name).Scan(&user.ID, &user.CreatedAt, &user.Version)
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

func (m UserModel) GetByID(id int) (*User, error) {

	query := `
		SELECT Id_users, Username, Email, Hashed_password, Avatar_path, Role, Birth_date, Created_at, Bio, Signature, Status, Version
		FROM users
		WHERE Id_users = ?;`

	var user User
	var birth sql.NullTime

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Avatar,
		&user.Role,
		&birth,
		&user.CreatedAt,
		&user.Bio,
		&user.Signature,
		&user.Status,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if birth.Valid {
		user.BirthDate = birth.Time
	}

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {

	query := `
		SELECT Id_users, Created_at, Username, Email, Hashed_password, Status, Version
		FROM users
		WHERE Email = ?;`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Status,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Activate(user *User) error {

	query := `
		UPDATE users
		SET Status = ?, Version = Version + 1
		WHERE Id_users = ? AND Version = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, UserStatus.Activated, user.ID, user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	user.Status = UserStatus.Activated

	return nil
}

func (m UserModel) Update(user *User) error {

	query := `
		UPDATE users
		SET Username = ?, Email = ?, Hashed_password = ?, Avatar_path = ?, Birth_date = ?, Bio = ?, Signature = ?, Updated_at = CURRENT_TIMESTAMP, Version = Version + 1
		WHERE Id_users = ? AND Version = ?;`

	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Avatar,
		user.BirthDate,
		user.Bio,
		user.Signature,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Email") {
					return ErrDuplicateEmail
				}
				if strings.Contains(mySQLError.Message, "Username") {
					return ErrDuplicateUsername
				}
			}
		default:
			return err
		}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT users.Id_users, users.Created_at, users.Username, users.Email, users.Hashed_password, users.Status, users.Version
		FROM users
		INNER JOIN tokens
		ON users.Id_users = tokens.Id_users
		WHERE tokens.Hash = ?
		AND tokens.Scope = ?
		AND tokens.Expiry > ?;`

	args := []any{hex.EncodeToString(tokenHash[:]), tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Status,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) DeleteExpired() error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE u
		FROM users u
		LEFT JOIN tokens t ON u.Id_users = t.Id_users
		WHERE u.Status = ? AND (t.Expiry IS NULL OR t.Expiry < CURRENT_TIMESTAMP);`

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, UserStatus.ToConfirm)
	if err != nil {
		return fmt.Errorf("failed to delete expired users: %w", err)
	}

	return nil
}
