package main

import (
	"Projet-Forum/internal/api"
	"Projet-Forum/internal/data"
	"Projet-Forum/internal/validator"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	var cfg config

	flag.BoolVar(&cfg.isHTTPS, "https", false, "use https")
	flag.IntVar(&cfg.port, "port", 4000, "HTTP service address")
	flag.StringVar(&cfg.apiURL, "api-url", "http://localhost:3000", "API URL")

	pemFilePath := flag.String("pem", "./pem/public.pem", "PEM file path")
	clientToken := flag.String("client-token", "", "Client token")

	dsn := flag.String("dsn", "", "MySQL DSN (data source name)")

	secret := flag.String("secret", "", "Secret API")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	addr := fmt.Sprintf(":%d", cfg.port)

	if secret == nil || *secret == "" {
		logger.Error("secret is required")
		os.Exit(1)
	}
	if pemFilePath == nil || *pemFilePath == "" {
		logger.Error("pem file path is required")
		os.Exit(1)
	}
	var pemKey []byte
	var err error
	if clientToken == nil || *clientToken == "" {
		clientToken, pemKey, err = getClientCredentials(cfg.apiURL, *secret, *pemFilePath)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	if pemKey == nil {
		pemKey, err = os.ReadFile(*pemFilePath)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		sessionManager: sessionManager,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		config:         &cfg,
		models:         data.NewModels(cfg.apiURL, *clientToken, pemKey),
	}

	server := http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),

		// timeouts setting, for security purposes. The server then automatically closes timed out connections
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	logger.Info("Starting server", slog.String("addr", server.Addr))

	if app.config.isHTTPS {
		// setting up the tls configuration
		// the CurvePreferences setting chosen here are the elliptic curves with assembly implementations
		// the MinVersion setting here specifies the minimum TLS version chosen (13 stands for 1.3 i.e. the last one at writing time)
		tlsConfig := &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
			MinVersion:       tls.VersionTLS13,
		}
		server.TLSConfig = tlsConfig

		// generate a signed certificate in tls folder for it to work
		// (for development use mkcert, for production, use let's encrypt)
		err = server.ListenAndServeTLS("./tls/localhost.pem", "./tls/localhost-key.pem")
	} else {

		// run the server through HTTP
		err = server.ListenAndServe()
	}

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func getClientCredentials(urlAPI, secret, pemFilePath string) (*string, []byte, error) {

	// creating connection to API
	apiClient := api.GetForClient(urlAPI, secret)

	// building request body with client credentials
	credentials := make(map[string]string)
	credentials["username"] = "Threadive Web"
	credentials["email"] = "web@threadive.com"
	v := validator.New()

	// getting client token
	clientToken, err := apiClient.GetClient(secret, credentials, v)
	if err != nil {
		return nil, nil, err
	}

	// fetching PEM public key
	var pem []byte
	if !fileExists(pemFilePath) {
		pem, err = apiClient.GetPEM(secret, pemFilePath, v)
		if err != nil {
			return nil, nil, err
		}
	} else {
		pem, err = os.ReadFile(pemFilePath)
		if err != nil {
			return nil, nil, err
		}
	}

	return clientToken, pem, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
