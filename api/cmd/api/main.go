package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/mailer"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

var (
	version = "v1"
)

type config struct {
	port int
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
		port     int
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
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "dsn", "", "MySQL Database DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "MySQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "MySQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "MySQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 50, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 100, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Threadive <no-reply@adebarbarin.com>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	secret := flag.String("secret", "", "API secret")

	frequency := flag.Duration("frequency", time.Hour*2, "expired tokens and unactivated users cleaning frequency")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Threadive API version:\t%s\n", version)
		os.Exit(0)
	}
	if cfg.smtp.username == "" || cfg.smtp.password == "" || cfg.smtp.host == "" {
		fmt.Println("SMTP credentials are required")
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(-4)}))

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

	app := &application{
		config:      cfg,
		logger:      logger,
		models:      data.NewModels(db),
		formDecoder: form.NewDecoder(),
		mailer:      mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	// TODO -> remove when leaving development phase
	if cfg.env == "development" && *secret != "" {
		user, err := app.models.Users.GetByEmail("api@api.com")
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				logger.Info("creating API user...")
				user = &data.User{
					Name:   "API",
					Email:  "api@api.com",
					Role:   data.UserRole.Secret,
					Status: data.UserStatus.HostSecret,
				}
				user.NoLogin()
				err = app.models.Users.Insert(user)
				if err != nil {
					logger.Error(err.Error())
					os.Exit(1)
				}
				hash := sha256.Sum256([]byte(*secret))
				token := data.Token{
					Plaintext: "",
					Hash:      hash[:],
					UserID:    user.ID,
					Expiry:    time.Now().Add(data.MaxDuration),
					Scope:     data.TokenScope.HostSecret,
				}
				err = app.models.Tokens.Insert(&token)
				if err != nil {
					logger.Error(err.Error())
					os.Exit(1)
				}
			} else {
				logger.Error(err.Error())
				os.Exit(1)
			}
		}
		if user == nil || user.ID < 1 {
			logger.Error("could not retrieve API user")
			os.Exit(1)
		}
		app.config.apiUserID = user.ID
	}
	*secret = ""
	// TODO <- END

	// Clean expired tokens every N duration
	go app.cleanExpiredTokens(*frequency, time.Hour*0)

	// Clean expired unactivated users every N duration
	go app.cleanExpiredUnactivatedUsers(*frequency, time.Hour)

	// Retrieving or generating RSA keys
	err = app.getPEM()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	credentials, err := app.encryptPEM([]byte(`{"email": "yellow@storm.com", "password": "Pa55word"}`))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	fmt.Printf("credentials: %s\n", hex.EncodeToString(credentials))

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
