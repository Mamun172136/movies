# Go Movies Backend API

A backend API for managing movies, built with Go, PostgreSQL, and JWT authentication.

## Features

- JWT authentication with refresh tokens
- CRUD operations for movies
- Genre management
- Movie poster fetching from TheMovieDB API
- PostgreSQL database
- RESTful API design

## Technologies

- Go 1.23.4
- PostgreSQL
- JWT for authentication
- Chi router
- Docker (for PostgreSQL)

## Setup

### Prerequisites

- Go 1.23.4 or later
- Docker (for running PostgreSQL)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Mamun172136/movies.git
   cd go-movies-backend
   docker-compose up -d
   go run ./cmd/api


# 🎬 Go Movies API Documentation
`http://localhost:8000/`
## 🔑 Authentication

### Login / Register User
`POST /authenticate`

**Request Body:**
```json
{
  "email": "admin1@example.com",
  "password":"test1234"
}
Success Response (200):

json
{
  "access_token": "eyJhbGciOiJIUz...",
  "refresh_token": "eyJhbGciOiJIUz..."
}
Sets HTTP-only cookie: refresh_token=eyJhbGci...

Refresh Access Token
GET /refresh

Requires Cookie:
refresh_token=eyJhbGci... (from login)

Success Response (200):

json
{
  "access_token": "eyJhbGciOiJIUz..."
}
Logout
GET /logout

Success Response (202):
(Empty body, clears refresh token cookie)

### Movies
- **GET /movies**: Retrieve a list of all movies.
Get All Movies
[
    {
        "id": 1,
        "title": "HighlandeR",
        "release_date": "0001-01-01T00:00:00Z",
        "runtime": 0,
        "mpaa_rating": "",
        "description": "",
        "image": "/8Z8dptJEypuLoOQro1WugD855YE.jpg"
    }
]

Success Response (200):
- **GET /movies/?id=1**: Retrieve details of a specific movie by ID.
### Genres

- **GET /genres**: Retrieve a list of all movie genres.
- **GET /movies/genres/?id=11**: Retrieve all movies belonging to a specific genre.

### Admin Routes (Requires Authentication)
- **Authorization: Bearer eyJhbGci...
- **GET /admin/movies**: Retrieve all movies (Admin only).
- **GET /admin/movies/?id=1**: Retrieve movie details for editing (Admin only).
- **PUT /admin/movie**: Insert a new movie (Admin only).
{
  "title": "Inception",
  "description": "A thief who steals corporate secrets...",
  "release_date": "2010-07-16T00:00:00Z",  // RFC3339 format
  "runtime": 148,
  "mpaa_rating": "PG-13"
 
}
- **PATCH /admin/movie**: Update movie details (Admin only).
{
    "ID": 1,
  "title": "HighlandeRR"
  
}
- **DELETE /admin/movies/?id=2**: Delete a movie (Admin only).
