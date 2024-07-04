package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/migrations"
	"database/sql"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/deatil/go-encoding/base62"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"golang.org/x/term"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	MaxDuration time.Duration = 1<<63 - 1
	TokenScope                = "host_secret"
)

type application struct {
	db         *sql.DB
	os         string
	migrations embed.FS
	reader     *bufio.Reader
	port       string
	secretAPI  string
	smtp       struct {
		host     string
		port     string
		username string
		password []byte
		sender   string
	}
	mysql struct {
		dsn   string
		host  string
		name  string
		admin struct {
			username string
			password []byte
		}
		api struct {
			username string
			password []byte
		}
		backend struct {
			username string
			password []byte
		}
	}
	backend struct {
		port string
	}
}

func main() {

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		fmt.Println("Caught signal:", s.String())
		fmt.Println("Shutting down...")

		os.Exit(0)
	}()

	var app application
	app.os = runtime.GOOS
	app.reader = bufio.NewReader(os.Stdin)
	app.migrations = migrations.Migrations
	var err error

	fmt.Println("###########################################")
	fmt.Printf("# Welcome to the Threadive migration tool!\n")
	fmt.Println("###########################################")
	fmt.Println()

	// #################################################################################
	// CONNECTING TO MySQL WITH ADMIN USER
	// #################################################################################

	fmt.Println("Let's connect to MySQL with an admin account (enough to create user, grant privileges and create a database)")
	fmt.Println()
	app.mysql.admin.username = app.readLine("Admin username:")
	fmt.Println()
	fmt.Print("Password: ")
	app.mysql.admin.password, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	// connect to MySQL
	app.db, err = app.openDB()
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer app.db.Close()

	// #################################################################################
	// SETTING THE NEW DATABASE CONNECTION
	// #################################################################################

	fmt.Println("###########################################")
	fmt.Println()
	fmt.Println("Let's set the new database:")
	fmt.Println()
	app.mysql.name = app.readLine("Database name:")
	fmt.Println()
	fmt.Println("Create a user to connect the API to the database:")
	fmt.Println()
	app.mysql.api.username = app.readLine("API username:")
	fmt.Println()
	fmt.Print("Generating password...")
	randomPassword := make([]byte, 12)
	_, err = rand.Read(randomPassword)
	if err != nil {
		fmt.Printf("Error generating password: %v\n", err)
		os.Exit(1)
	}
	app.mysql.api.password = base62.StdEncoding.Encode(randomPassword)
	fmt.Println()

	fmt.Println("Creating empty API database...")
	err = app.createDB()
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	fmt.Println("Creating MySQL user for API...")
	err = app.createUser()
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	// #################################################################################
	// SETTING THE .ENV VARIABLES AND FILE
	// #################################################################################

	fmt.Println("###########################################")
	fmt.Println()
	fmt.Println("Let's set the SMTP configuration:")
	fmt.Println()
	app.smtp.host = app.readLine("SMTP Host:")
	fmt.Println()
	app.smtp.port = app.readLine("SMTP Port:")
	fmt.Println()
	app.smtp.sender = app.readLine("SMTP Sender:")
	fmt.Println()
	app.smtp.username = app.readLine("SMTP Username:")
	fmt.Println()
	fmt.Print("SMTP Password: ")
	app.smtp.password, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	fmt.Println("###########################################")
	fmt.Println()
	fmt.Println("Let's set the API:")
	fmt.Println()
	app.port = app.readLine("API port:")
	fmt.Println()

	fmt.Println("Generating the API secret token...")
	app.generateSecret()
	fmt.Println()

	fmt.Println("Writing configurations to environment file...")

	switch app.os {
	case "windows":
		app.genEnvFile()
	case "linux":
		app.genEnvrcFile()
	default:
		fmt.Println("Unsupported OS")
		os.Exit(1)
	}
	fmt.Println()

	fmt.Println("###########################################")
	fmt.Println()
	fmt.Println("Let's set the Backend:")
	fmt.Println()
	app.backend.port = app.readLine("Backend port:")
	fmt.Println()
	fmt.Println("Create a user to connect the Backend to the database: (to access the sessions Table)")
	fmt.Println()
	app.mysql.backend.username = app.readLine("Backend username:")
	fmt.Println()
	fmt.Print("Generating password...")
	randomPassword = make([]byte, 12)
	_, err = rand.Read(randomPassword)
	if err != nil {
		fmt.Printf("Error generating password: %v\n", err)
		os.Exit(1)
	}
	app.mysql.backend.password = base62.StdEncoding.Encode(randomPassword)
	fmt.Println()

	// FIXME -> create backend user!!!

	fmt.Println("Writing configurations to environment file...")

	var env strings.Builder
	port := fmt.Sprintf("PORT=%s\n", app.backend.port)
	apiURL := fmt.Sprintf("API_URL=\"http://localhost:%s\"\n", app.port)
	secretToken := fmt.Sprintf("API_HOST_SECRET=\"%s\"\n", app.secretAPI)
	dbDSN := fmt.Sprintf("DB_DSN=\"%s:%s@/%s?parseTime=true\"", app.mysql.backend.username, string(app.mysql.backend.password), app.mysql.name)

	switch app.os {
	case "windows":
		env.WriteString(port)
		env.WriteString(apiURL)
		env.WriteString(secretToken)
		env.WriteString(dbDSN)
		err := os.WriteFile("backend.env", []byte(env.String()), 0644)
		if err != nil {
			fmt.Println("Error creating backend.env file")
			os.Exit(1)
		}

		fmt.Println("backend.env file created successfully")
	case "linux":
		env.WriteString(fmt.Sprintf("export %s", port))
		env.WriteString(fmt.Sprintf("export %s", apiURL))
		env.WriteString(fmt.Sprintf("export %s", secretToken))
		env.WriteString(fmt.Sprintf("export %s", dbDSN))
		err := os.WriteFile("backend.envrc", []byte(env.String()), 0644)
		if err != nil {
			fmt.Println("Error creating backend.envrc file")
			os.Exit(1)
		}
	default:
		fmt.Println("Unsupported OS")
		os.Exit(1)
	}
	fmt.Println()

	// #################################################################################
	// RUNNING THE MIGRATIONS
	// #################################################################################

	fmt.Println("###########################################")
	fmt.Println()
	fmt.Println("Setting up the database...")
	fmt.Println()

	d, err := iofs.New(app.migrations, ".")
	if err != nil {
		fmt.Println(err)
		return
	}

	app.mysql.dsn = fmt.Sprintf("%s:%s@/%s?parseTime=true", app.mysql.admin.username, string(app.mysql.admin.password), app.mysql.name)

	dbMigrations, err := migrate.NewWithSourceInstance("iofs", d, fmt.Sprintf("mysql://%s", app.mysql.dsn))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dbMigrations.Up()
	if err != nil && err.Error() != "no change" {
		fmt.Println(err)
		return
	}

	err = app.insertSecretToken()
	if err != nil {
		fmt.Printf("Error inserting secret token: %v\n", err)
	}

	fmt.Println("API database successfully created!")
}

func (app *application) genEnvFile() {
	var env strings.Builder
	env.WriteString("DB_HOST=\"localhost\"\n")
	env.WriteString("DB_PORT=\"3306\"\n")
	env.WriteString(fmt.Sprintf("DB_USER=\"%s\"\n", app.mysql.api.username))
	env.WriteString(fmt.Sprintf("DB_PASSWORD=\"%s\"\n", string(app.mysql.api.password)))
	env.WriteString(fmt.Sprintf("DB_DATABASE=\"%s\"\n", app.mysql.name))
	env.WriteString("DB_ARG=\"parseTime=true\"\n")
	env.WriteString("DB_NETWORK=\"tcp\"\n")
	env.WriteString(fmt.Sprintf("PORT=%s\n", app.port))
	env.WriteString(fmt.Sprintf("DB_DSN=\"%s:%s@/%s?parseTime=true\"\n", app.mysql.api.username, string(app.mysql.api.password), app.mysql.name))
	env.WriteString(fmt.Sprintf("SMTP_SENDER=\"%s\"\n", app.smtp.sender))
	env.WriteString(fmt.Sprintf("SMTP_USERNAME=\"%s\"\n", app.smtp.username))
	env.WriteString(fmt.Sprintf("SMTP_PASS=\"%s\"\n", string(app.smtp.password)))
	env.WriteString(fmt.Sprintf("SMTP_HOST=\"%s\"\n", app.smtp.host))
	env.WriteString(fmt.Sprintf("SMTP_PORT=%s\n", app.smtp.port))
	env.WriteString(fmt.Sprintf("API_HOST_SECRET=\"%s\"\n", app.secretAPI))

	err := os.WriteFile(".env", []byte(env.String()), 0644)
	if err != nil {
		fmt.Println("Error creating .env file")
		os.Exit(1)
	}

	fmt.Println(".env file created successfully")
}

func (app *application) genEnvrcFile() {
	var env strings.Builder
	env.WriteString("export DB_HOST=\"localhost\"\n")
	env.WriteString("export DB_PORT=\"3306\"\n")
	env.WriteString(fmt.Sprintf("export DB_USER=\"%s\"\n", app.mysql.api.username))
	env.WriteString(fmt.Sprintf("export DB_PASSWORD=\"%s\"\n", string(app.mysql.api.password)))
	env.WriteString(fmt.Sprintf("export DB_DATABASE=\"%s\"\n", app.mysql.name))
	env.WriteString("export DB_ARG=\"parseTime=true\"\n")
	env.WriteString("export DB_NETWORK=\"tcp\"\n")
	env.WriteString(fmt.Sprintf("export PORT=%s\n", app.port))
	env.WriteString(fmt.Sprintf("export DB_DSN=\"%s:%s@/%s?parseTime=true\"\n", app.mysql.api.username, string(app.mysql.api.password), app.mysql.name))
	env.WriteString(fmt.Sprintf("export SMTP_SENDER=\"%s\"\n", app.smtp.sender))
	env.WriteString(fmt.Sprintf("export SMTP_USERNAME=\"%s\"\n", app.smtp.username))
	env.WriteString(fmt.Sprintf("export SMTP_PASS=\"%s\"\n", string(app.smtp.password)))
	env.WriteString(fmt.Sprintf("export SMTP_HOST=\"%s\"\n", app.smtp.host))
	env.WriteString(fmt.Sprintf("export SMTP_PORT=%s\n", app.smtp.port))
	env.WriteString(fmt.Sprintf("export API_HOST_SECRET=\"%s\"\n", app.secretAPI))

	err := os.WriteFile(".envrc", []byte(env.String()), 0644)
	if err != nil {
		fmt.Println("Error creating .envrc file")
		os.Exit(1)
	}

	fmt.Println(".envrc file created successfully")
}

func (app *application) readLine(prompt string) string {

	fmt.Println(prompt)
	fmt.Print(">>> ")

	input, err := app.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}

	input = strings.TrimSpace(input)

	return input
}

func (app *application) generateSecret() {

	randomBytes := make([]byte, 64)

	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}

	app.secretAPI = base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
}
func (app *application) openDB() (*sql.DB, error) {

	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", app.mysql.admin.username, string(app.mysql.admin.password), app.mysql.name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (app *application) createDB() error {

	query := fmt.Sprintf(`
		CREATE DATABASE IF NOT EXISTS %s;`, app.mysql.name)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := app.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return err
}

func (app *application) createUser() error {

	query := fmt.Sprintf(`
		CREATE USER '%s'@'localhost' IDENTIFIED BY '%s';`, app.mysql.api.username, string(app.mysql.api.password))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := app.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		GRANT ALL PRIVILEGES ON %s . * TO '%s'@'localhost';`, app.mysql.name, app.mysql.api.username)

	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (app *application) createBackendUser() error {

	query := fmt.Sprintf(`
		CREATE USER '%s'@'localhost'
    	IDENTIFIED BY '%s';`, app.mysql.backend.username, string(app.mysql.backend.password))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := app.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		GRANT SELECT, INSERT, UPDATE, DELETE
    	ON %s.sessions
    	TO '%s'@'localhost';`, app.mysql.name, app.mysql.backend.username)

	stmt, err = tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (app *application) insertSecretToken() error {

	hash := sha256.Sum256([]byte(app.secretAPI))
	token := hash[:]

	expiry := time.Now().Add(MaxDuration)

	query := fmt.Sprintf(`
		INSERT INTO %s.tokens (Hash, Id_users, Expiry, Scope)
		VALUES (?, ?, ?, ?);`, app.mysql.name)

	args := []any{hex.EncodeToString(token), 2, expiry, TokenScope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := app.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}
