CREATE TABLE IF NOT EXISTS genres (
    genre_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS dramas (
    drama_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_year INT,
    main_actors VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS genres_ans_dramas (
    id SERIAL PRIMARY KEY,
    genre_id INT,
    drama_id INT,
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id),
    FOREIGN KEY (drama_id) REFERENCES dramas(movie_id)
);