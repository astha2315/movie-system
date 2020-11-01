package models

type Movie struct {
	MovieId    int     `db:"id" json:"id"`
	MovieName  string  `db:"name" json:"name"`
	Popularity float64 `db:"popularity" json:"99popularity"`
	ImdbScore  float64 `db:"imdb_score" json:"imdb_score"`
	Director   string  `db:"director" json:"director"`
	Status     int     `db:"status" json:"status"`
}

type Genre struct {
	GenreId   int    `db:"id" json:"id"`
	GenreName string `db:"name" json:"name"`
	Status    int    `db:"status" json:"status"`
}

type MovieGenre struct {
	MovieGenreId int `db:"id" json:"id"`
	GenreId      int `db:"genre_id" json:"genre_id"`
	MovieId      int `db:"movie_id" json:"movie_id"`
	Status       int `db:"status" json:"status"`
}

type MovieJSON struct {
	MovieId    int      `db:"id" json:"id"`
	MovieName  string   `db:"name" json:"name"`
	Genre      []string `json:"genre"`
	Popularity float64  `db:"rating" json:"99popularity"`
	ImdbScore  float64  `db:"imdb_score" json:"imdb_score"`
	Director   string   `db:"director" json:"director"`
}

type MovieObj struct {
	MovieJson []MovieJSON ` json:"movies"`
}
