package main

import (
	"encoding/json"
	"fmt"
	"log"
	"movie-system/controllers"
	"movie-system/db"
	"movie-system/utils"
	"net/http"

	"github.com/casbin/casbin"
)

type Enforcer struct {
	enforcer casbin.Enforcer
}

func main() {
	fmt.Println("Hello")

	_, err := db.DBConnect()
	if err != nil {
		panic(err)
	}

	cE := casbin.NewEnforcer("model.conf", "policy.csv")

	enforcer := Enforcer{enforcer: *cE}

	http.HandleFunc("/movies/list", enf(&enforcer, controllers.GetMovies))
	http.HandleFunc("/movies/add", enf(&enforcer, controllers.AddMovies))
	http.HandleFunc("/movies/genres/list", enf(&enforcer, controllers.GetAllGenres))
	http.HandleFunc("/movies/remove", enf(&enforcer, controllers.RemoveMovies))
	http.HandleFunc("/movies/list/genre", enf(&enforcer, controllers.GetMoviesByGenre))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func enf(e *Enforcer, next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Header.Get("Authorization")
		fmt.Println(user)

		method := r.Method
		path := r.URL.Path

		result := e.enforcer.Enforce(user, path, method)
		if !result {

			json.NewEncoder(w).Encode(utils.ResponseJson{http.StatusForbidden, nil, 0, "Method not allowed"})
			return
		}

		next(w, r)

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Movie System!")

}
