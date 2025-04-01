package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type PostgresEDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second*3

func (m *PostgresEDBRepo) AllMovies() ([]*models.Movie, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select
	id, title, release_date, runtime, mpaa_rating, description, coalesce(image, '' ),
	created_at, update_at
	from movies
	order by title
	
	`
	rows, err:= m.DB.QueryContext(ctx, query)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie

	for rows.Next(){
		var movie models.Movie
		err := rows.Scan(&movie.ID,
			 &movie.Title,
			  &movie.ReleaseDate,
			  &movie.Runtime,
			  &movie.MPAARating,
			  &movie.Description,
			  &movie.Image,
			  &movie.CreatedAt,
			  &movie.UpdatedAt,
		)
		if err !=nil{
			return nil, err
		}

		movies = append(movies, &movie)
	}
	return movies, nil
}

func (m *PostgresEDBRepo) Connection() *sql.DB{
	return m.DB
}

func (m *PostgresEDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresEDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}