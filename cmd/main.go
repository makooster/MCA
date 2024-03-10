package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/makooster/MCA/pkg/model"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func (app *application) run() {
	r := mux.NewRouter()

	r.HandleFunc("/", app.HomeHandler)
	r.HandleFunc("/user", app.UserHandler)
	r.HandleFunc("/dramas", app.DramaHandler).Methods("GET")
	r.HandleFunc("/dramas", app.createDramaHandler).Methods("POST")
	r.HandleFunc("/dramas/{id}", app.getDramaHandler).Methods("GET")
	r.HandleFunc("/dramas/{id}", app.updateDramaHandler).Methods("PUT")
	r.HandleFunc("/dramas/{id}", app.deleteDramaHandler).Methods("DELETE")
	r.HandleFunc("/actors", app.ActorHandler).Methods("GET")
	r.HandleFunc("/actors", app.createActorHandler).Methods("POST")
	r.HandleFunc("/actors/{id}", app.getActorHandler).Methods("GET")
	r.HandleFunc("/actors/{id}", app.updateActorHandler).Methods("PUT")
	r.HandleFunc("/actors/{id}", app.deleteActorHandler).Methods("DELETE")
	r.HandleFunc("/dramas", app.getDramasHandler).Methods("GET")
	r.HandleFunc("/actors", app.getActorsHandler).Methods("GET")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1234Asdf@localhost:5432/go?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
