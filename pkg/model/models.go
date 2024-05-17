package model

import (
	"database/sql"
	"log"
	"os"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Doramas DoramaModel
	Actors ActorModel
	Genres GenreModel
	Users UserModel
	Tokens TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Doramas: DoramaModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Actors: ActorModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Genres: GenreModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Permissions: PermissionModel{DB: db},
		Tokens: TokenModel{DB: db}, 
		Users: UserModel{DB: db},
	}
}
