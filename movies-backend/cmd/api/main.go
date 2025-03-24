package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const port = 8000

type application struct {
	DSN string
	DB repository.DatabaseRepo
}

func main() {
	// set application config
	var app application

	//read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()
	//connect to database
	conn , err:= app.connectToDB()
	if err != nil{
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresEDBRepo{DB: conn}
	defer conn.Close()
	
	log.Println("Starting application on Port", port)
	// http.HandleFunc("/", Hello)
	//start web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err !=nil{
		log.Fatal(err)
	}
}