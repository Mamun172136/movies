package main

import (
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
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
	// var movies []models.Movie
	// rd, _ := time.Parse("2006-01-02", "1986-03-07")
	// highlander := models.Movie{
	// 	ID:1,
	// 	Title: "Highlander",
	// 	ReleaseDate:rd ,
	// 	MPAARating: "R",
	// 	Runtime: 116,
	// 	Description: "nice movie",
	// 	CreatedAt:time.Now() ,
	// 	UpdatedAt: time.Now(),
	// }

	// movies = append(movies, highlander)

	movies, err :=  app.DB.AllMovies()
	if err != nil{
		app.errorJSON(w, err)
		fmt.Println(err)
	}

	// out, err := json.Marshal(movies)
	// if err != nil{
	// 	fmt.Println(err)
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(out)

	_ =  app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)  
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID: user.ID,
				FirstName: user.FirstName,
				LastName: user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) MovieCatalog(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

// func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {


// 	id := chi.URLParam(r, "id")
	
// 	// Method 2: Direct context access
// 	rctx := chi.RouteContext(r.Context())
// 	if rctx != nil {
// 		log.Printf("Handler route params: %+v", rctx.URLParams)
// 	}

// 	log.Printf("Handling movie ID: %s", id)
	
// 	if id == "" {
// 		app.errorJSON(w, errors.New("movie ID not found in URL parameters"))
// 		return
// 	}



// 	// id := chi.URLParam(r, "id")
// 	fmt.Printf("Received ID: '%s'\n", id) // Add this line to debug
// 	log.Printf("Request path: %s", r.URL.Path)
// 	log.Printf("----Handler context: %+v", chi.RouteContext(r.Context()))
// 	movieID, err := strconv.Atoi(id)
// 	if err != nil {
// 		app.errorJSON(w, err)
// 		return
// 	}

// 	movie, err := app.DB.OneMovie(movieID)
// 	if err != nil {
// 		app.errorJSON(w, err)
// 		return
// 	}

// 	_ = app.writeJSON(w, http.StatusOK, movie)
// }

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	// First try to get ID from URL path (e.g., /movies/1)
    id := chi.URLParam(r, "id")
    
    // If empty, check query params (e.g., /movies?id=1)
    if id == "" {
        id = r.URL.Query().Get("id")
    }
    
    if id == "" {
        app.errorJSON(w, errors.New("movie ID is required"))
        return
    }

	

	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movie)
}

func (app *application) MovieForEdit(w http.ResponseWriter, r *http.Request) {
// First try to get ID from URL path (e.g., /movies/1)
id := chi.URLParam(r, "id")
    
// If empty, check query params (e.g., /movies?id=1)
if id == "" {
	id = r.URL.Query().Get("id")
}

if id == "" {
	app.errorJSON(w, errors.New("movie ID is required"))
	return
}
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		movie,
		genres,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.DB.AllGenres()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, genres)
}

func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := app.readJSON(w, r, &movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// try to get an image
	movie = app.getPoster(movie)

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()


	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// now handle genres
	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse {
		Error: false,
		Message: "movie updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) getPoster(movie models.Movie) models.Movie {
	type TheMovieDB struct {
		Page int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}

	client := &http.Client{}
	theUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)

	// https://api.themoviedb.org/3/search/movie?api_key=b41447e6319d1cd467306735632ba733&query=Die+Hard

	req, err := http.NewRequest("GET", theUrl+"&query="+url.QueryEscape(movie.Title), nil)
	if err != nil {
		log.Println(err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}

	var responseObject TheMovieDB

	json.Unmarshal(bodyBytes, &responseObject)

	if len(responseObject.Results) > 0 {
		movie.Image = responseObject.Results[0].PosterPath
	}

	return movie
}

func (app *application) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var payload models.Movie

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie.Title = payload.Title
	movie.ReleaseDate = payload.ReleaseDate
	movie.Description = payload.Description
	movie.MPAARating = payload.MPAARating
	movie.RunTime = payload.RunTime
	movie.UpdatedAt = time.Now()

	err = app.DB.UpdateMovie(*movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.UpdateMovieGenres(movie.ID, payload.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error: false,
		Message: "movie updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	// First try to get ID from URL path (e.g., /movies/1)
id := chi.URLParam(r, "id")
    
// If empty, check query params (e.g., /movies?id=1)
if id == "" {
	id = r.URL.Query().Get("id")
}

if id == "" {
	app.errorJSON(w, errors.New("movie ID is required"))
	return
}
	
	mid, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteMovie(mid)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error: false,
		Message: "movie deleted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) AllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	// First try to get ID from URL path (e.g., /movies/1)
id := chi.URLParam(r, "id")
    
// If empty, check query params (e.g., /movies?id=1)
if id == "" {
	id = r.URL.Query().Get("id")
}

if id == "" {
	app.errorJSON(w, errors.New("movie ID is required"))
	return
}
log.Println("GENRE ID:", id)
	
	mid, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movies, err := app.DB.AllMovies(mid)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, movies)
}

