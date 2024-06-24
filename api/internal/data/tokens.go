package data

import (
	"ForumAPI/internal/validator"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

const (
	MaxDuration time.Duration = 1<<63 - 1
)

type tokenScope struct {
	Activation     string
	Authentication string
	Refresh        string
	Client         string
	HostSecret     string
}

var (
	TokenScope = &tokenScope{
		Activation:     "activation",
		Authentication: "authentication",
		Refresh:        "refresh",
		Client:         "client",
		HostSecret:     "host_secret",
	}
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int       `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userID int, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 64)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 86, "token", "must be 86 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	if errors.Is(err, ErrDuplicateToken) {
		token, err = generateToken(userID, ttl, scope)
		if err != nil {
			return nil, err
		}

		err = m.Insert(token)
	}
	return token, err
}

func (m TokenModel) Insert(token *Token) error {

	query := `
		INSERT INTO tokens (Hash, Id_users, Expiry, Scope)
		VALUES (?, ?, ?, ?);`

	args := []any{hex.EncodeToString(token.Hash), token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		switch {
		case errors.As(err, &mySQLError):
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "Hash") {
					return ErrDuplicateToken
				}
			} else {
				return err
			}
		default:
			return err
		}
	}

	return nil
}

func (m TokenModel) DeleteAllForUser(scope string, userID int) error {

	query := `
		DELETE FROM tokens
       	WHERE Scope = ? AND Id_users = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if scope == "*" {
		query = `
			DELETE FROM tokens
       		WHERE Id_users = ?;`

		_, err := m.DB.ExecContext(ctx, query, userID)
		return err
	}

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}

func (m TokenModel) DeleteExpired() error {

	query := `
		DELETE FROM tokens
		WHERE Expiry < CURRENT_TIMESTAMP;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query)
	return err
}
