package repository

import (
	"context"
	"project-app-bioskop/internal/data/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CinemaRepoInterface interface {
	GetAllCinemas(ctx context.Context, limit, offset int) ([]entity.Cinema, error)
	GetCinemaByID(ctx context.Context, id int) (entity.Cinema, error)
	GetStudiosByCinemaID(ctx context.Context, cinemaID int) ([]entity.Studio, error)
	CountCinemas(ctx context.Context) (int, error)
}

type CinemaRepo struct {
	DB *pgxpool.Pool
}

func NewCinemaRepo(db *pgxpool.Pool) CinemaRepoInterface {
	return &CinemaRepo{DB: db}
}

// GetAllCinemas retrieves all cinemas with pagination
func (r *CinemaRepo) GetAllCinemas(ctx context.Context, limit, offset int) ([]entity.Cinema, error) {
	query := `SELECT id, name, location, created_at FROM cinemas ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cinemas []entity.Cinema
	for rows.Next() {
		var c entity.Cinema
		if err := rows.Scan(&c.ID, &c.Name, &c.Location, &c.CreatedAt); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, c)
	}
	return cinemas, nil
}

// GetCinemaByID retrieves cinema by ID
func (r *CinemaRepo) GetCinemaByID(ctx context.Context, id int) (entity.Cinema, error) {
	query := `SELECT id, name, location, created_at FROM cinemas WHERE id = $1`
	var c entity.Cinema
	err := r.DB.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Location, &c.CreatedAt)
	if err != nil {
		return c, err
	}
	return c, nil
}

// GetStudiosByCinemaID retrieves all studios for a cinema
func (r *CinemaRepo) GetStudiosByCinemaID(ctx context.Context, cinemaID int) ([]entity.Studio, error) {
	query := `SELECT id, cinema_id, name, total_seats FROM studios WHERE cinema_id = $1 ORDER BY id`
	rows, err := r.DB.Query(ctx, query, cinemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studios []entity.Studio
	for rows.Next() {
		var s entity.Studio
		if err := rows.Scan(&s.ID, &s.CinemaID, &s.Name, &s.TotalSeats); err != nil {
			return nil, err
		}
		studios = append(studios, s)
	}
	return studios, nil
}

// CountCinemas returns total cinema count
func (r *CinemaRepo) CountCinemas(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM cinemas`
	var count int
	err := r.DB.QueryRow(ctx, query).Scan(&count)
	return count, err
}
