package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
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

	addr := flag.String("addr", ":4000", "HTTP service address")

	// if you want to enable a MySQL database (mainly for the sessions management)
	dsn := flag.String("dsn", "forum:forum-Pa55w0rd/serverTemplate?parseTime=true", "MySQL DSN (data source name)")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// if you want a MySQL database linked to your web server (mainly for the sessions management)
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
	// if you want to store the sessionIDs in a MySQL database, with db being the database pool
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		sessionManager: sessionManager,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
	}

	// setting up the tls configuration
	// the CurvePreferences setting chosen here are the elliptic curves with assembly implementations
	// the MinVersion setting here specifies the minimum TLS version chosen (13 stands for 1.3 i.e. the last one at writing time)
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	server := http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,

		// timeouts setting, for security purposes. The server then automatically closes timed out connections
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	logger.Info("Starting server", slog.String("addr", server.Addr))

	// generate a signed certificate in tls folder for it to work
	// (for development use mkcert, for production, use let's encrypt)
	err = server.ListenAndServeTLS("./tls/localhost.pem", "./tls/localhost-key.pem")
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
