package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"log"
)

type Actor struct {
	ActorId int    `json:"actor_id"`
	Name    string `json:"name"`
	FilmID  string `json:"film_id"`
}

type ActorModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (am *ActorModel) Get(id int) (*Actor, error) {
	query := `
        SELECT actor_id, full_name, film_id
        FROM actors
        WHERE actor_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	actor := &Actor{}
	err := am.DB.QueryRowContext(ctx, query, id).Scan(&actor.ActorId, &actor.Name)
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
		INSERT INTO actors (full_name, film_id) 
		VALUES ($1, $2) 
		RETURNING actor_id
		`
	args := []interface{}{actor.Name}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return am.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ActorId)
}

func (am *ActorModel) Update(actor *Actor) error {
	query := `
	 UPDATE actors
	 SET full_name = $1, film_id = $2
	 WHERE actor_id = $3
	 RETURNING actor_id
	`
	args := []interface{}{actor.Name, actor.ActorId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return am.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ActorId)
}

func (am *ActorModel) Delete(id int) error {
	query := `
	 DELETE FROM actors 
	 WHERE actor_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := am.DB.ExecContext(ctx, query, id)
	return err
}
