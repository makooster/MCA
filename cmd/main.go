package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"context"
	"fmt"
	"os"
	"time"
	_ "github.com/lib/pq"
	"github.com/makooster/MCA/pkg/model"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host string
		port int
		username string
		password string
		sender string
	}
}

type application struct {
	config config
	logger log.Logger
	models model.Models
}

func main() {
	// Declare an instance of the config struct.
	var cfg config
	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://postgres:qwerty@localhost:5432/mca?sslmode=disable", "PostgreSQL DSN")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "0f1d85c09e6d8e", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "e89654b1c53c45", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Printf("database connection pool established")

	app := &application {
		config: cfg,
		logger: *logger,
		models: model.NewModels(db),
	}
	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr: 			fmt.Sprintf(":%d", cfg.port),
		Handler: 		app.routes(),
		IdleTimeout:	time.Minute,
		ReadTimeout: 	10 * time.Second,
		WriteTimeout: 	30 * time.Second,
	}
	
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
	
