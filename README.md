# Movie Catalog App

This is a simple movie catalog app that allows users to manage movies, genres, and actors. The app provides RESTful APIs for basic CRUD operations on movies, genres, and actors.

## Table of Contents
- [Project members](#project-members)
- [API Endpoints](#api-endpoints)
- [Database Structure](#database-structure)

## Project Members
|Full name|ID|
|---|---|
|Magzhan Akhmadi|22B030517|
|Sarkyt Asylai|22B030585|
|Davlatova Altyn|22B030334|
|Birlikzhanova Aruzhan|22B030329|

## API Endpoints
- **GET /movies**: Retrieve all movies.
- **GET /movies/{id}**: Retrieve a specific movie by ID.
- **POST /movies**: Create a new movie.
- **PUT /movies/{id}**: Update a specific movie.
- **DELETE /movies/{id}**: Delete a specific movie.
- ...

- **GET /genres**: Retrieve all genres.
- **GET /genres/{id}**: Retrieve a specific genre by ID.
- **POST /genres**: Create a new genre.
- **PUT /genres/{id}**: Update a specific genre.
- **DELETE /genres/{id}**: Delete a specific genre.
- ...

- **GET /actors**: Retrieve all actors.
- **GET /actors/{id}**: Retrieve a specific actor by ID.
- **POST /actors**: Create a new actor.
- **PUT /actors/{id}**: Update a specific actor.
- **DELETE /actors/{id}**: Delete a specific actor.
- ...



## Database Structure

Here's a simplified representation of the database structure:

```sql
Database Structure:
Table: genres

Columns:
genre_id: Auto-incremented identifier (primary key).
name: Text field for the genre name.
Table: movies

Columns:
movie_id: Auto-incremented identifier (primary key).
title: Text field for the movie title.
description: Text field for the movie description.
release_year: Integer field for the movie release year.
main_actors: Text field for the main actors.
Table: actors

Columns:
actor_id: Auto-incremented identifier (primary key).
name: Text field for the actor's name.
