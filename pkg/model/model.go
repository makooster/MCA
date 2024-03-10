package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Genres GenreModel
	Movies DramaModel
	Actors ActorModel
}

// type GenreModel struct {
// 	DB       *sql.DB
// 	InfoLog  *log.Logger
// 	ErrorLog *log.Logger
// }

type DramaModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type ActorModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (dm *DramaModel) GetDramas() []Drama {
	return dramas
}

func (am *ActorModel) GetActors() []Actor {
	return actors
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Genres: GenreModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Movies: DramaModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Actors: ActorModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
