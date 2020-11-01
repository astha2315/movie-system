package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "moviedb"
)

func DBConnect() (*DB, error) {

	// db, err := sqlx.Open("postgres", "host=localhost user=postgres password=password dbname=moviedb sslmode=disable")

	// db, err := sqlx.Open("postgres", "user=postgres  dbname=moviedb sslmode=disable")

	conString := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + "disable"
	db, err := sqlx.Open("postgres", conString)

	if err != nil {
		log.Fatalln(err)
	}
	return &DB{db}, err
}

func SqlxConnect() (*DB, error) {
	// return db
	db, err := DBConnect()

	return db, err
}
