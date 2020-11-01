package dao

import (
	"fmt"
	"movie-system/db"
	"movie-system/models"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type movieDao struct{}

type MovieDaoIF interface {
	GetAllMovies() ([]*models.Movie, error)
	GetMoviesByGenre(genres []string) ([]*models.Movie, error)
	GetMoviesByName(movieName string) ([]*models.Movie, error)

	RemoveMovie(movieId int) (int, error)
	GetAllGenres() ([]*models.Genre, error)
	AddMovieList(movieList []*models.Movie) (int, error)
	AddGenreList(genreList []*models.Genre) (int, error)
	GetMoviesByRange(limit int) ([]*models.Movie, error)
	GetGenresByRange(limit int) ([]*models.Genre, error)

	AddMovieGenreList(movieGenreList []*models.MovieGenre) (int, int, error)
}

func MovieDao() MovieDaoIF {
	return &movieDao{}
}

func (self *movieDao) GetAllMovies() ([]*models.Movie, error) {

	var movies []*models.Movie
	db, ConnectionErrs := db.SqlxConnect()
	if ConnectionErrs != nil {
		return movies, ConnectionErrs
	}

	sqlStatement := `select  name,popularity,director,imdb_score from movie where status=0 order by id`
	err := db.Select(&movies, sqlStatement)
	// var err error
	if err != nil {
		fmt.Println(err)
		return movies, err
	}
	defer db.Close()
	return movies, err
}

func (self *movieDao) GetMoviesByGenre(genres []string) ([]*models.Movie, error) {

	var movies []*models.Movie
	db, ConnectionErrs := db.SqlxConnect()
	if ConnectionErrs != nil {
		return movies, ConnectionErrs
	}

	sqlStatement := `select m.name from movie m 
	join movie_genre mg on m.id=mg.movie_id and m.status=0
	join genre g on mg.genre_id=g.id and mg.status=0
	where g.name in(?) and g.status=0`

	query, args, err := sqlx.In(sqlStatement, genres)

	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&movies, query, args...)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	return movies, err
}

func (self *movieDao) GetMoviesByName(movieName string) ([]*models.Movie, error) {

	var movies []*models.Movie
	db, ConnectionErrs := db.SqlxConnect()
	if ConnectionErrs != nil {
		return movies, ConnectionErrs
	}

	sqlStatement := `select movie_name from movie order by id where name=$1 and status=0`
	err := db.Select(&movies, sqlStatement, movieName)
	// var err error
	if err != nil {
		fmt.Println(err)
		return movies, err
	}
	defer db.Close()
	return movies, err
}

func (self *movieDao) RemoveMovie(movieId int) (int, error) {
	db, ConnectionErrs := db.DBConnect()

	if ConnectionErrs != nil {
		fmt.Println(ConnectionErrs)
		return movieId, ConnectionErrs
	}

	sqlStatement := `UPDATE movie SET status=1 WHERE id=$1 and status=0`
	_, custMatQueryErr := db.Exec(sqlStatement, movieId)

	if custMatQueryErr != nil {
		fmt.Println("the error is", custMatQueryErr)
		return movieId, custMatQueryErr
	}

	defer db.Close()
	return movieId, custMatQueryErr
}

func (self *movieDao) GetAllGenres() ([]*models.Genre, error) {

	var genres []*models.Genre
	db, ConnectionErrs := db.SqlxConnect()
	if ConnectionErrs != nil {
		return genres, ConnectionErrs
	}

	sqlStatement := `select name,id from genre order by id`
	err := db.Select(&genres, sqlStatement)
	// var err error
	if err != nil {
		fmt.Println(err)
		return genres, err
	}
	defer db.Close()
	return genres, err
}

func ReplaceSQL(stmt, pattern string, len int) string {
	pattern += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(pattern, len))
	n := 0
	for strings.IndexByte(stmt, '?') != -1 {
		n++
		param := "$" + strconv.Itoa(n)
		stmt = strings.Replace(stmt, "?", param, 1)
	}
	return strings.TrimSuffix(stmt, ",")
}

func (self *movieDao) AddMovieList(movieList []*models.Movie) (int, error) {

	db, errs := db.SqlxConnect()
	if errs != nil {
		return 0, errs
	}
	vals := []interface{}{}
	for _, row := range movieList {
		vals = append(vals, row.MovieName, row.Popularity, row.ImdbScore, row.Director, row.Status)
	}
	sqlStr := `INSERT INTO movie (name,popularity,imdb_score,director,status) VALUES %s`
	sqlStr = ReplaceSQL(sqlStr, "(?,?,?,?,?)", len(movieList))
	sqlStr += "RETURNING id"
	stmt, prepErr := db.Prepare(sqlStr)
	if prepErr != nil {
		fmt.Printf("error in preparing sql stmt", prepErr)
		return 0, prepErr
	}
	fmt.Println("****")
	fmt.Println(stmt)
	fmt.Println("****")
	res, excErr := stmt.Exec(vals...)
	if excErr != nil {
		fmt.Printf("error in executing sql str", excErr)
		return 0, excErr
	}
	rowsEffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	// lastInsertedId, err := res.LastInsertId()
	// if err != nil {
	// 	return 0, 0, err
	// }
	defer db.Close()
	return int(rowsEffected), excErr
}

func (self *movieDao) AddGenreList(genreList []*models.Genre) (int, error) {

	db, errs := db.SqlxConnect()
	if errs != nil {
		return 0, errs
	}
	vals := []interface{}{}
	for _, row := range genreList {
		vals = append(vals, row.GenreName, row.Status)
	}
	sqlStr := `INSERT INTO genre (name,status) VALUES %s`
	sqlStr = ReplaceSQL(sqlStr, "(?,?)", len(genreList))
	stmt, prepErr := db.Prepare(sqlStr)
	if prepErr != nil {
		fmt.Printf("error in preparing sql stmt", prepErr)
		return 0, prepErr
	}
	fmt.Println("****")
	fmt.Println(stmt)
	fmt.Println("****")
	res, excErr := stmt.Exec(vals...)
	if excErr != nil {
		fmt.Printf("error in executing sql str", excErr)
		return 0, excErr
	}
	rowsEffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	defer db.Close()
	return int(rowsEffected), excErr
}

func (self *movieDao) AddMovieGenreList(movieGenreList []*models.MovieGenre) (int, int, error) {

	db, errs := db.SqlxConnect()
	if errs != nil {
		return 0, 0, errs
	}
	vals := []interface{}{}
	for _, row := range movieGenreList {
		vals = append(vals, row.MovieId, row.GenreId, row.Status)
	}
	sqlStr := `INSERT INTO movie_genre (movie_id,genre_id,status) VALUES %s`
	sqlStr = ReplaceSQL(sqlStr, "(?,?,?)", len(movieGenreList))
	stmt, prepErr := db.Prepare(sqlStr)
	if prepErr != nil {
		fmt.Printf("error in preparing sql stmt", prepErr)
		return 0, 0, prepErr
	}
	fmt.Println("****")
	fmt.Println(stmt)
	fmt.Println("****")
	res, excErr := stmt.Exec(vals...)
	if excErr != nil {
		fmt.Printf("error in executing sql str", excErr)
		return 0, 0, excErr
	}
	rowsEffected, err := res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	defer db.Close()
	return int(rowsEffected), int(lastInsertedId), excErr
}

func (self *movieDao) GetMoviesByRange(limit int) ([]*models.Movie, error) {

	var movies []*models.Movie
	db, ConnectionErrs := db.SqlxConnect()
	if ConnectionErrs != nil {
		return movies, ConnectionErrs
	}

	sqlStatement := `select id,name from movie order by id desc limit $1 `
	err := db.Select(&movies, sqlStatement, limit)
	// var err error
	if err != nil {
		fmt.Println(err)
		return movies, err
	}
	defer db.Close()
	return movies, err
}

func (self *movieDao) GetGenresByRange(limit int) ([]*models.Genre, error) {

	var genres []*models.Genre
	db, ConnectionErrs := db.SqlxConnect()
	if ConnectionErrs != nil {
		return genres, ConnectionErrs
	}

	sqlStatement := `select id,name from genre order by id desc limit $1 `
	err := db.Select(&genres, sqlStatement, limit)
	// var err error
	if err != nil {
		fmt.Println(err)
		return genres, err
	}
	defer db.Close()
	return genres, err
}
