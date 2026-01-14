package usecase

import (
	"context"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
)

type MovieUseCaseInterface interface {
	GetAllMovies(ctx context.Context, page, limit int) ([]dto.MovieResponse, dto.Pagination, error)
	GetMovieByID(ctx context.Context, id int) (dto.MovieResponse, error)
}

type MovieUseCase struct {
	Repo *repository.Repository
}

func NewMovieUseCase(repo *repository.Repository) MovieUseCaseInterface {
	return &MovieUseCase{Repo: repo}
}

// GetAllMovies retrieves all movies with pagination
func (u *MovieUseCase) GetAllMovies(ctx context.Context, page, limit int) ([]dto.MovieResponse, dto.Pagination, error) {
	movies, pagination, err := u.Repo.Movie.GetAllMovies(ctx, page, limit)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	var response []dto.MovieResponse
	for _, m := range movies {
		response = append(response, dto.MovieResponse{
			ID:              m.ID,
			Title:           m.Title,
			PosterURL:       m.PosterURL,
			Genres:          m.Genres,
			Rating:          m.Rating,
			DurationMinutes: m.DurationMinutes,
		})
	}

	return response, pagination, nil
}

// GetMovieByID retrieves a movie by ID
func (u *MovieUseCase) GetMovieByID(ctx context.Context, id int) (dto.MovieResponse, error) {
	movie, err := u.Repo.Movie.GetMovieByID(ctx, id)
	if err != nil {
		return dto.MovieResponse{}, err
	}

	return dto.MovieResponse{
		ID:              movie.ID,
		Title:           movie.Title,
		PosterURL:       movie.PosterURL,
		Genres:          movie.Genres,
		Rating:          movie.Rating,
		DurationMinutes: movie.DurationMinutes,
	}, nil
}
