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

// type application struct {
// 	config config
// 	models model.Models
// }

type application struct {
	config config
	logger log.Logger
	models model.Models
}


func (app *application) run() {
	r := mux.NewRouter()
	//Home 
	r.HandleFunc("/", app.HomeHandler)
	r.HandleFunc("/user", app.UserHandler)
	

	//Healthcheck 
	r.HandleFunc("/check", app.healthcheckHandler).Methods("GET")
	
	// r.HandleFunc("/doramas", app.getDoramasHandler).Methods("GET")
	// r.HandleFunc("/actors", app.getActorsHandler).Methods("GET")

	//GET methods

	// r.HandleFunc("/genres", app.GenreHandler).Methods("GET")
	r.HandleFunc("/list", app.getDoramaListHandler).Methods("GET")
	r.HandleFunc("/doramas/{id}", app.getDoramaHandler).Methods("GET")
	r.HandleFunc("/actors", app.getActorsHandler).Methods("GET")
	r.HandleFunc("/actors/{id}", app.getActorHandler).Methods("GET")

	//POST methods
	r.HandleFunc("/doramas", app.createDoramaHandler).Methods("POST")
	r.HandleFunc("/actors", app.createActorHandler).Methods("POST")

	//DELETE methods 
	r.HandleFunc("/doramas/{id}", app.deleteDoramaHandler).Methods("DELETE")
	r.HandleFunc("/actors/{id}", app.deleteActorHandler).Methods("DELETE")

	//PUT methods 
	r.HandleFunc("/doramas/{id}", app.updateDoramaHandler).Methods("PUT")
	r.HandleFunc("/actors/{id}", app.updateActorHandler).Methods("PUT")
	

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8080", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://postgres:qwerty@localhost:5432/mca?sslmode=disable", "PostgreSQL DSN")

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
