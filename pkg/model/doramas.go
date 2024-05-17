package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"log"
	"fmt"
)

type Dorama struct {
	DoramaId    int    `json:"dorama_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Duration    int    `json:"duration"`
	MainActors  string `json:"main_actors"`
	GenreId     int    `json:"genre_id"`
}

type DoramaModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m DoramaModel) GetAll(title string, releaseYear int,filters Filters) ([]*Dorama, Metadata, error) {
	// Retrieve all doramas from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), dorama_id, title, description, release_year, duration, main_actors
		FROM doramas
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (release_year = $2 OR $2 = 1)
		ORDER BY %s %s, dorama_id
		LIMIT $3 OFFSET $4`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our placeholder parameter values in a slice.
	args := []interface{}{title, releaseYear,filters.limit(), filters.offset()}

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
	totalRecords := 0

	var doramas []*Dorama
	for rows.Next() {
		var dorama Dorama
		err := rows.Scan(&totalRecords, &dorama.DoramaId, &dorama.Title, &dorama.Description, &dorama.ReleaseYear, &dorama.Duration, &dorama.MainActors)
		if err != nil {
			return nil, Metadata{}, err
		}
		doramas = append(doramas, &dorama)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate Metadata struct based on total record count and pagination parameters.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// Return doramas and metadata.
	return doramas, metadata, nil
}

func (dm *DoramaModel) Get(id int) (*Dorama, error) {
	query := `
	SELECT dorama_id, title, description, release_year, duration, main_actors, genre_id
	FROM doramas
	WHERE dorama_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dorama := &Dorama{}
	err := dm.DB.QueryRowContext(ctx, query, id).Scan(&dorama.DoramaId, &dorama.Title, &dorama.Description, &dorama.ReleaseYear,&dorama.Duration, &dorama.MainActors, &dorama.GenreId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("film not found")
		} else {
			return nil, err
		}
	}

	return dorama, nil
}

func (dm *DoramaModel) Insert(dorama *Dorama) error {
	query := `
		INSERT INTO doramas (title, description, release_year, duration, main_actors, genre_id) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		where dorama_id = $7
		RETURNING dorama_id,title, description, release_year, duration, main_actors, genre_id
		`
	args := []interface{}{dorama.Title, dorama.Description, dorama.ReleaseYear, dorama.Duration ,dorama.MainActors, dorama.GenreId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return dm.DB.QueryRowContext(ctx, query, args...).Scan(&dorama.DoramaId)
}

func (dm *DoramaModel) Update(dorama *Dorama) error {
    query := `
        UPDATE doramas
        SET title = $1, description = $2, release_year = $3, duration = $4, main_actors = $5, genre_id = $6
        WHERE dorama_id = $7
        RETURNING dorama_id
    `
    args := []interface{}{dorama.Title,dorama.Description,dorama.ReleaseYear,dorama.Duration,dorama.MainActors,dorama.GenreId,dorama.DoramaId}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
	
    return dm.DB.QueryRowContext(ctx, query, args...).Scan(&dorama.DoramaId)
}

func (dm *DoramaModel) Delete(id int) error {
	query := `
        DELETE FROM doramas
        WHERE dorama_id = $1
        `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := dm.DB.ExecContext(ctx, query, id)
	return err
}

