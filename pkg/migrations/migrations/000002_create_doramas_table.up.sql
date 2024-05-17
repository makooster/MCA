CREATE TABLE IF NOT EXISTS doramas (
    dorama_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_year INT,
    duration INT,
    main_actors text,
    genre_id integer,
    foreign key (genre_id) references genres(genre_id)
);