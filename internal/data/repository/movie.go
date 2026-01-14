package repository

import (
	"context"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepoInterface interface {
	GetAllMovies(ctx context.Context, page, limit int) ([]entity.Movie, dto.Pagination, error)
	GetMovieByID(ctx context.Context, id int) (entity.Movie, error)
}

type MovieRepo struct {
	DB *pgxpool.Pool
}

func NewMovieRepo(db *pgxpool.Pool) MovieRepoInterface {
	return &MovieRepo{DB: db}
}

// GetAllMovies retrieves all movies with pagination
func (r *MovieRepo) GetAllMovies(ctx context.Context, page, limit int) ([]entity.Movie, dto.Pagination, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM movies`
	if err := r.DB.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, dto.Pagination{}, err
	}

	// Get paginated movies
	query := `SELECT id, title, poster_url, genres, rating, review_count, 
			  release_date, duration_in_minutes, release_status, created_at, updated_at
			  FROM movies 
			  ORDER BY created_at DESC
			  LIMIT $1 OFFSET $2`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, dto.Pagination{}, err
	}
	defer rows.Close()

	var movies []entity.Movie
	for rows.Next() {
		var m entity.Movie
		if err := rows.Scan(
			&m.ID, &m.Title, &m.PosterURL, &m.Genres, &m.Rating, &m.ReviewCount,
			&m.ReleaseDate, &m.DurationMinutes, &m.ReleaseStatus, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, dto.Pagination{}, err
		}
		movies = append(movies, m)
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		TotalPages:   (total + limit - 1) / limit,
		TotalRecords: total,
		Limit:        limit,
	}

	return movies, pagination, nil
}

// GetMovieByID retrieves a movie by ID
func (r *MovieRepo) GetMovieByID(ctx context.Context, id int) (entity.Movie, error) {
	query := `SELECT id, title, poster_url, genres, rating, review_count, 
			  release_date, duration_in_minutes, release_status, created_at, updated_at
			  FROM movies WHERE id = $1`

	var m entity.Movie
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Title, &m.PosterURL, &m.Genres, &m.Rating, &m.ReviewCount,
		&m.ReleaseDate, &m.DurationMinutes, &m.ReleaseStatus, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return m, err
	}
	return m, nil
}
