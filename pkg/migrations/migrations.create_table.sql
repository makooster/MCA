CREATE TABLE genres (
    genre_id SERIAL PRIMARY KEY,
    genre_name VARCHAR(255) NOT NULL
);

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

create table if not exists actors (
    id serial primary key,
    full_name varchar(255),
    dorama_id integer,
    foreign key (dorama_id) references doramas(dorama_id)
 );


