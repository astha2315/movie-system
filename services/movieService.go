package services

import (
	"fmt"
	"movie-system/dao"
	"movie-system/models"
)

type movieService struct{}

type MovieServiceIF interface {
	GetMovies() ([]*models.Movie, error)
	AddMovies(movieObj *models.MovieObj) error
	GetAllGenres() ([]*models.Genre, error)
	RemoveMovies(id int) (int, error)
	GetMoviesByGenre(genres []string) ([]*models.Movie, error)
}

func MovieService() MovieServiceIF {
	return &movieService{}
}

func (self *movieService) GetMovies() ([]*models.Movie, error) {

	movieData, err := dao.MovieDao().GetAllMovies()

	if err != nil {
		return nil, err
	}
	return movieData, nil

}

func (self *movieService) AddMovies(movieObj *models.MovieObj) error {

	//insert all movies
	//get the existing genres
	// make a list of all genres to be inserted
	//check if the new genre to be inserted already exists
	//if yes, keep a tab on its id
	//for the new genres not exisiting, insert them in the db
	//now we have the genre ids as well
	// as well as the movie ids
	//now map accordingly (movieId,genreId) to be inserted in movie_genre

	existingGenres, err := dao.MovieDao().GetAllGenres()

	if err != nil {
		return err
	}

	movieData := movieObj.MovieJson

	//parsing the input data

	var newGenres []string
	var newMovies []*models.Movie

	if movieData != nil && len(movieData) > 0 {
		for _, m := range movieData {
			if len(m.Genre) > 0 {
				newGenres = append(newGenres, m.Genre...)
			}

			mov := &models.Movie{
				MovieName:  m.MovieName,
				Popularity: m.Popularity,
				ImdbScore:  m.ImdbScore,
				Director:   m.Director,
				Status:     0,
			}

			newMovies = append(newMovies, mov)
		}
	}

	movieRowsAffected, err := dao.MovieDao().AddMovieList(newMovies)
	fmt.Println(movieRowsAffected)

	if err != nil {
		return err
	}

	exisitingGenresMap := make(map[string]int)

	if existingGenres != nil && len(existingGenres) > 0 {
		for _, genre := range existingGenres {

			if _, ok := exisitingGenresMap[genre.GenreName]; !ok {
				exisitingGenresMap[genre.GenreName] = genre.GenreId
			}
		}
	}

	newGenresMap := make(map[string]int)
	var newGenreList []*models.Genre

	if len(newGenres) > 0 {
		for _, newGenre := range newGenres {

			//checking if the new genre already exist in the exisitng genre map

			if _, ok := exisitingGenresMap[newGenre]; !ok {

				// now we keep the new genre in the new genre map, because the genres can repeat,
				//  making 0 as the id of the genre
				if _, ok := newGenresMap[newGenre]; !ok {

					newGenresMap[newGenre] = 0

					g := new(models.Genre)
					g.GenreName = newGenre
					g.Status = 0

					newGenreList = append(newGenreList, g)
				}
			}
		}
	}

	//inserting the list of genres

	genreRowsAffected, err := dao.MovieDao().AddGenreList(newGenreList)
	fmt.Println(genreRowsAffected)

	if err != nil {
		return err
	}

	//now in order to insert into movie_genre, we need the data of movies and genre
	insertedMovies, err := dao.MovieDao().GetMoviesByRange(movieRowsAffected)

	if err != nil {
		return err
	}

	insertedGenres, err := dao.MovieDao().GetGenresByRange(genreRowsAffected)

	if err != nil {
		return err
	}

	// making the map of the inserted movies and genres
	insertedMovieMap := make(map[string]int)

	if insertedMovies != nil && len(insertedMovies) > 0 {
		for _, m := range insertedMovies {

			if _, ok := insertedMovieMap[m.MovieName]; !ok {
				insertedMovieMap[m.MovieName] = m.MovieId
			}
		}
	}

	insertedGenreMap := make(map[string]int)

	if insertedGenres != nil && len(insertedGenres) > 0 {
		for _, g := range insertedGenres {

			if _, ok := insertedGenreMap[g.GenreName]; !ok {
				insertedGenreMap[g.GenreName] = g.GenreId
			}
		}
	}

	// making the structure for movie_genre table

	var movieGenreList []*models.MovieGenre

	for _, m := range movieData {

		var movieId int
		movieId = insertedMovieMap[m.MovieName]

		for _, g := range m.Genre {

			mg := new(models.MovieGenre)
			mg.MovieId = movieId

			var genreId int

			if id, ok := insertedGenreMap[g]; ok {
				genreId = id
			} else {
				genreId = exisitingGenresMap[g]
			}

			mg.GenreId = genreId

			movieGenreList = append(movieGenreList, mg)

		}
	}

	movieGenreRowsAffected, lastMovieGenreId, err := dao.MovieDao().AddMovieGenreList(movieGenreList)
	fmt.Println(movieGenreRowsAffected)
	fmt.Println(lastMovieGenreId)

	if err != nil {
		return err
	}

	return nil

}

func (self *movieService) GetAllGenres() ([]*models.Genre, error) {

	genreData, err := dao.MovieDao().GetAllGenres()

	if err != nil {
		return nil, err
	}
	return genreData, nil

}

func (self *movieService) RemoveMovies(id int) (int, error) {

	idRemoved, err := dao.MovieDao().RemoveMovie(id)

	if err != nil {
		return 0, err
	}
	return idRemoved, nil

}

func (self *movieService) GetMoviesByGenre(genres []string) ([]*models.Movie, error) {

	movies, err := dao.MovieDao().GetMoviesByGenre(genres)

	if err != nil {
		return nil, err
	}
	return movies, nil
}
