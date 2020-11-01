package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"movie-system/models"
	"movie-system/services"
	"movie-system/utils"
	"net/http"
	"strconv"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {

	fmt.Println("In get api for movies")

	services := services.MovieService()
	movies, err := services.GetMovies()

	if err != nil {
		json.NewEncoder(w).Encode(utils.ResponseJson{0, nil, 0, err.Error()})
	}

	// return c.JSON(http.StatusOK, response)
	json.NewEncoder(w).Encode(movies)

}

func RemoveMovies(w http.ResponseWriter, r *http.Request) {

	fmt.Println("In remove api for movies")

	id, _ := r.URL.Query()["id"]

	movieId, err := strconv.Atoi(id[0])
	if err != nil {
		json.NewEncoder(w).Encode(utils.ResponseJson{0, nil, 0, err.Error()})
	}
	services := services.MovieService()
	movies, err := services.RemoveMovies(movieId)

	if err != nil {
		json.NewEncoder(w).Encode(utils.ResponseJson{0, nil, 0, err.Error()})
	}

	// return c.JSON(http.StatusOK, response)
	json.NewEncoder(w).Encode(movies)

}

func AddMovies(w http.ResponseWriter, r *http.Request) {

	fmt.Println("In add api for movies")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "%+v", string(reqBody))

	// var movieJSON []*models.MovieJSON
	movieJSON := new(models.MovieObj)

	if err := json.Unmarshal(reqBody, &movieJSON); err != nil {
		fmt.Println(err)
		panic(err)
	}

	services := services.MovieService()
	_ = services.AddMovies(movieJSON)

}

func GetAllGenres(w http.ResponseWriter, r *http.Request) {

	fmt.Println("In get api for movies")

	services := services.MovieService()
	genre, err := services.GetAllGenres()

	if err != nil {
		json.NewEncoder(w).Encode(utils.ResponseJson{0, nil, 0, err.Error()})
	}

	// return c.JSON(http.StatusOK, response)
	json.NewEncoder(w).Encode(genre)

}

func GetMoviesByGenre(w http.ResponseWriter, r *http.Request) {

	services := services.MovieService()

	genres, _ := r.URL.Query()["genre"]
	movies, err := services.GetMoviesByGenre(genres)

	if err != nil {
		json.NewEncoder(w).Encode(utils.ResponseJson{0, nil, 0, err.Error()})
	}

	// return c.JSON(http.StatusOK, response)
	json.NewEncoder(w).Encode(movies)

}
