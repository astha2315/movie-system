CREATE TABLE movie(
id serial PRIMARY KEY,
name VARCHAR(50) NOT NULL,
imdb_score NUMERIC(10,2) NOT NULL,
popularity NUMERIC(10,2) NOT NULL ,
director VARCHAR(50) NOT NULL,
status SMALLINT DEFAULT 0 NOT NULL
);


CREATE TABLE genre(
id serial PRIMARY KEY,
name VARCHAR(50) NOT NULL,
status SMALLINT DEFAULT 0 NOT NULL
);


CREATE TABLE movie_genre(
id serial ,
movie_id integer REFERENCES movie(id) NOT NULL,
genre_id integer REFERENCES genre(id) NOT NULL,
PRIMARY KEY(movie_id,genre_id),
status SMALLINT DEFAULT 0 NOT NULL

);
