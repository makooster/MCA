package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"log"
	"fmt"
)

type Actor struct {
	ActorId int    `json:"id"`
	Name    string `json:"full_name"`
	DoramaID  int  `json:"dorama_id"`
}

type ActorModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}



func (m ActorModel) GetAll(fullName string, doramaID int, filters Filters) ([]*Actor, Metadata, error) {
	// Retrieve all actors from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, full_name, dorama_id
		FROM actors
		WHERE (to_tsvector('simple', full_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (dorama_id = $2 OR $2 = 1)
		ORDER BY %s %s, dorama_id
		LIMIT $3 OFFSET $4`,
	filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our placeholder parameter values in a slice.
	args := []interface{}{fullName, doramaID, filters.PageSize, filters.Page}

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

	var actors []*Actor
	for rows.Next() {
		var actor Actor
		err := rows.Scan(&totalRecords, &actor.ActorId, &actor.Name, &actor.DoramaID)
		if err != nil {
			return nil, Metadata{}, err
		}
		actors = append(actors, &actor)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate Metadata struct based on total record count and pagination parameters.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// Return actors and metadata.
	return actors, metadata, nil
}
func (am *ActorModel) Get(id int) (*Actor, error) {
	query := `
        SELECT id, full_name, dorama_id
        FROM actors
        WHERE id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	actor := &Actor{}
	err := am.DB.QueryRowContext(ctx, query, id).Scan(&actor.ActorId, &actor.Name, &actor.DoramaID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("actor not found")
		} else {
			return nil, err
		} 
	}

	return actor, nil
}

func (am *ActorModel) Insert(actor *Actor) error {
	query := `
		INSERT INTO actors (full_name, dorama_id) 
		VALUES ($1, $2) 
		RETURNING id
		`
	args := []interface{}{actor.Name, actor.DoramaID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return am.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ActorId)
}

func (am *ActorModel) Update(actor *Actor) error {
    query := `
        UPDATE actors
        SET full_name = $1, dorama_id = $2
        WHERE id = $3
        RETURNING id
    `
    args := []interface{}{actor.Name, actor.DoramaID, actor.ActorId}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return am.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ActorId)
}

func (am *ActorModel) Delete(id int) error {
	query := `
	 DELETE FROM actors 
	 WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := am.DB.ExecContext(ctx, query, id)
	return err
} 
