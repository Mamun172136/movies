package main

import (
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (app *application)Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct{
		Status string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status: "active",
		Message: "Go movies up and running",
		Version: "1.0.0",
	}

	out, err := json.Marshal(payload)
	if err != nil{
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request){
	var movies []models.Movie
	rd, _ := time.Parse("2006-01-02", "1986-03-07")
	highlander := models.Movie{
		ID:1,
		Title: "Highlander",
		ReleaseDate:rd ,
		MPAARating: "R",
		Runtime: 116,
		Description: "nice movie",
		CreatedAt:time.Now() ,
		UpdatedAt: time.Now(),
	}

	movies = append(movies, highlander)

	out, err := json.Marshal(movies)
	if err != nil{
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}