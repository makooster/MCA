package model

import (
	"context"
	"database/sql"
	"log"
	"time"
	"fmt"
	"errors"
)

type Genre struct {
	GenreID      int    `json:"genre_id"`
	GenreName    string `json:"genre_name"`
}

type GenreModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m GenreModel) GetAll(GenreName string, GenreID int, filters Filters) ([]*Genre, Metadata, error) {
    // Construct the SQL query
    query := fmt.Sprintf(
        `
        SELECT count(*) OVER(), genre_id, genre_name
        FROM genres
        WHERE (to_tsvector('simple', genre_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (genre_id = $2 OR $2 = 1)
        ORDER BY %s %s, genre_id
        LIMIT $3 OFFSET $4`,
        filters.sortColumn(), filters.sortDirection())

    // Create a context with a 3-second timeout.
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Organize our placeholder parameter values in a slice.
    args := []interface{}{GenreName, GenreID, filters.PageSize, filters.Page}

    // Use QueryContext to execute the query.
    rows, err := m.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, Metadata{}, err
    }
    defer func() {
        if err := rows.Close(); err != nil {
            m.ErrorLog.Println(err)
        }
    }()

    // Declare a totalRecords variable
    var totalRecords int

    var genres []*Genre
    for rows.Next() {
        var genre Genre
        err := rows.Scan(&totalRecords, &genre.GenreID, &genre.GenreName)
        if err != nil {
            return nil, Metadata{}, err
        }
        genres = append(genres, &genre)
    }

    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    // Generate Metadata struct based on total record count and pagination parameters.
    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

    // Return genres and metadata.
    return genres, metadata, nil
}


func (gm *GenreModel) Get(id int) (*Genre,error) {

	query := `
		SELECT genre_id, genre_name
		FROM genres
		where genre_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	genre:= &Genre{}
	err := gm.DB.QueryRowContext(ctx, query, id).Scan(&genre.GenreID, &genre.GenreName)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("genre not found")
		} else {
			return nil, err
		} 
	}

	return genre, nil
}
 
func (gm *GenreModel) Insert(genre *Genre) error {
	query := `
  INSERT INTO genres (genre_name) 
  VALUES ($1) 
  RETURNING genre_id
 `
	args := []interface{}{genre.GenreID, genre.GenreName}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return gm.DB.QueryRowContext(ctx, query, args...).Scan(&genre.GenreID)
}

func (gm *GenreModel) Update(genre *Genre) error {
	query := `
  UPDATE genres
  SET genre_name = $1
  WHERE genre_id = $2
  RETURNING genre_id
 `
	args := []interface{}{genre.GenreName, genre.GenreID}
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
