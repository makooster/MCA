package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Genre struct {
	GenreID int    `json:"genre_id"`
	Name    string `json:"name"`
}

var genres = []Genre{
	{
		GenreID: 1,
		Name:    "Исторические",
	},
	{
		GenreID: 2,
		Name:    "Комедия",
	},
	{
		GenreID: 3,
		Name:    "Приключения",
	},
	{
		GenreID: 4,
		Name:    "Мелодрама",
	},
	{
		GenreID: 5,
		Name:    "Триллеры",
	},
}

func GetGenres() []Genre {
	return genres
}

type GenreModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (gm *GenreModel) Get(id int) (*Genre, error) {
	for _, g := range genres {
		if g.GenreID == id {
			return &g, nil
		}
	}
	return nil, errors.New("Genre not found")
}

func (gm *GenreModel) Insert(genre *Genre) error {
	query := `
  INSERT INTO genres (genre_id, name) 
  VALUES ($1, $2) 
  RETURNING genre_id
 `
	args := []interface{}{genre.GenreID, genre.Name}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return gm.DB.QueryRowContext(ctx, query, args...).Scan(&genre.GenreID)
}

func (gm *GenreModel) Update(genre *Genre) error {
	query := `
  UPDATE genres
  SET name = $1
  WHERE genre_id = $2
  RETURNING genre_id
 `
	args := []interface{}{genre.Name, genre.GenreID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return gm.DB.QueryRowContext(ctx, query, args...).Scan(&genre.GenreID)
}

func (gm *GenreModel) Delete(id int) error {
	query := `
  DELETE FROM genres
  WHERE genre_id = $1
 `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := gm.DB.ExecContext(ctx, query, id)
	return err
}
