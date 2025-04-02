package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8000

type application struct {
	DSN string
	DB repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey	string
}

func main() {
	// set application config
	var app application

	//read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.Parse()
	//connect to database
	conn , err:= app.connectToDB()
	if err != nil{
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresEDBRepo{DB: conn}
	defer conn.Close()

	app.auth = Auth{
		Issuer: app.JWTIssuer,
		Audience: app.JWTAudience,
		Secret: app.JWTSecret,
		TokenExpiry: time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath: "/",
		CookieName: "__Host-refresh_token",
		CookieDomain: app.CookieDomain,
	}
	log.Println("Starting application on Port", port)
	// http.HandleFunc("/", Hello)
	//start web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err !=nil{
		log.Fatal(err)
	}
}