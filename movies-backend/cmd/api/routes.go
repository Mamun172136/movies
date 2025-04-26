package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.StripSlashes) // Add this
	mux.Use(app.enableCORS)

	// Add debug middleware AFTER routes
    mux.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // This will now show proper params
            if rctx := chi.RouteContext(r.Context()); rctx != nil {
                log.Printf("Route params: %+v", rctx.URLParams)
            }
            next.ServeHTTP(w, r)
        })
    })

	mux.Get("/",app.Home)
	mux.Get("/movies", app.AllMovies)

	mux.Get("/movie", app.GetMovie)
	mux.Post("/authenticate", app.authenticate)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)
	mux.Get("/genres", app.AllGenres)
	mux.Get("/movies/genres", app.AllMoviesByGenre)


	mux.Route("/admin", func(mux chi.Router){
		mux.Use(app.authRequired)
		mux.Get("/movie", app.MovieForEdit)  // QUERY
		mux.Put("/insert/movie", app.InsertMovie)
		mux.Get("/movies", app.MovieCatalog)
		mux.Patch("/movie", app.UpdateMovie)
		mux.Delete("/movie", app.DeleteMovie)
	})
	return mux
}