package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/mailer"
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	version = "v1"
)

type config struct {
	port int64
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int64
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
	pem struct {
		privateKey []byte
		publicKey  []byte
	}
	apiUserID int
}

type application struct {
	config      config
	logger      *slog.Logger
	models      data.Models
	formDecoder *form.Decoder
	mailer      mailer.Mailer
	wg          sync.WaitGroup
}

func main() {
	// Creating the config struct and retreiving the flags
	var cfg config

	flag.Int64Var(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "dsn", "", "MySQL Database DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "MySQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "MySQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "MySQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 50, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 100, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "", "SMTP host")
	flag.Int64Var(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	frequency := flag.Duration("frequency", time.Hour*2, "expired tokens and unactivated users cleaning frequency")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	// parsing the flags
	flag.Parse()

	// displaying the version and exit (if flag version found)
	if *displayVersion {
		fmt.Printf("Threadive API current version:\t%s\n", version)
		os.Exit(0)
	}

	// checking the SMTP info & retreive info from .env file if OS => Windows
	if cfg.smtp.username == "" || cfg.smtp.password == "" || cfg.smtp.host == "" {
		if runtime.GOOS == "windows" {
			err := cfg.loadEnv()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading environment variables: %s\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("SMTP credentials are required")
			os.Exit(1)
		}
	}

	// creating the logger with level corresponding to the environment (development|staging|production)
	var logger *slog.Logger
	if cfg.env == "development" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	// opening the database connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	// creating the application
	app := &application{
		config:      cfg,
		logger:      logger,
		models:      data.NewModels(db),
		formDecoder: form.NewDecoder(),
		mailer:      mailer.New(cfg.smtp.host, int(cfg.smtp.port), cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	// Clean expired tokens every N duration with no timeout
	go app.cleanExpiredTokens(*frequency, time.Hour*0)

	// Clean expired unactivated users every N duration with 1 hour timeout
	go app.cleanExpiredUnactivatedUsers(*frequency, time.Hour)

	// Retrieving or generating RSA keys
	err = app.getPEM()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Running the server
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {

	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (cfg *config) loadEnv() error {

	err := godotenv.Load()
	if err != nil {
		return err
	}

	cfg.port, err = strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		return err
	}
	cfg.db.dsn = os.Getenv("DB_DSN")
	cfg.smtp.sender = os.Getenv("SMTP_SENDER")
	cfg.smtp.username = os.Getenv("SMTP_USERNAME")
	cfg.smtp.password = os.Getenv("SMTP_PASS")
	cfg.smtp.host = os.Getenv("SMTP_HOST")
	cfg.smtp.port, err = strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 64)
	if err != nil {
		return err
	}

	return nil
}
