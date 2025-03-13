package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 8000

type application struct {
}

func main() {
	// set application config
	var app application

	//read from command line

	//connect to database
	log.Println("Starting application on Port", port)
	http.HandleFunc("/", Hello)
	//start web server
	err:= http.ListenAndServe(fmt.Sprintf("%d", port), app.routes())
	if err !=nil{
		log.Fatal(err)
	}
}