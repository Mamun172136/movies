# Go Project

This is a Go project that can be compiled and executed locally using the `go build` command. It does not require Docker or installing Go globally on your system.

## Prerequisites

Before you begin, make sure that you have Go installed on your machine.

1. **Install Go**:  
   Download and install Go from the official website: [Go Downloads](https://golang.org/dl/)

2. **Verify Go installation**:  
   After installing Go, verify the installation by running the following command in your terminal and make sure that you have **Go 1.23.4** (or a compatible version) installed on your machine.:

   ```bash
   go version
   git clone <https://github.com/Mamun172136/movies.gitt>
   cd <project-directory>/movies-backend
   go get -u
   go run ./cmd/api


# Movies API

## API Endpoints

### Authentication

- **POST /authenticate**: Authenticate a user and get a token.
- **POST /signup**: Register a new user.
- **POST /login**: Authenticate a user and get a token.
- **GET /refresh**: Refresh the authentication token.
- **GET /logout**: Logout the user.

### Movies

- **GET /movies**: Retrieve a list of all movies.
- **GET /movies/{id}**: Retrieve details of a specific movie by ID.
- **POST /movies**: Add a new movie (Admin only).
- **PUT /movies/{id}**: Update details of a specific movie by ID (Admin only).
- **DELETE /movies/{id}**: Delete a specific movie by ID (Admin only).

### Genres

- **GET /genres**: Retrieve a list of all movie genres.
- **GET /movies/genres/{id}**: Retrieve all movies belonging to a specific genre.

### Admin Routes (Requires Authentication)

- **GET /admin/movies**: Retrieve all movies (Admin only).
- **GET /admin/movies/{id}**: Retrieve movie details for editing (Admin only).
- **PUT /admin/movies/0**: Insert a new movie (Admin only).
- **PATCH /admin/movies/{id}**: Update movie details (Admin only).
- **DELETE /admin/movies/{id}**: Delete a movie (Admin only).
