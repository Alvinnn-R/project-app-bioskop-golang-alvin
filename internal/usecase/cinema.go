package usecase

import (
	"context"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
)

type CinemaUseCaseInterface interface {
	GetAllCinemas(ctx context.Context, page, limit int) ([]dto.CinemaResponse, dto.Pagination, error)
	GetCinemaByID(ctx context.Context, id int) (dto.CinemaResponse, error)
}

type CinemaUseCase struct {
	Repo *repository.Repository
}

func NewCinemaUseCase(repo *repository.Repository) CinemaUseCaseInterface {
	return &CinemaUseCase{Repo: repo}
}

// GetAllCinemas retrieves all cinemas with pagination
func (u *CinemaUseCase) GetAllCinemas(ctx context.Context, page, limit int) ([]dto.CinemaResponse, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	cinemas, err := u.Repo.Cinema.GetAllCinemas(ctx, limit, offset)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	totalRecords, err := u.Repo.Cinema.CountCinemas(ctx)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	totalPages := totalRecords / limit
	if totalRecords%limit != 0 {
		totalPages++
	}

	var response []dto.CinemaResponse
	for _, c := range cinemas {
		response = append(response, dto.CinemaResponse{
			ID:        c.ID,
			Name:      c.Name,
			Location:  c.Location,
			CreatedAt: c.CreatedAt,
		})
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}

	return response, pagination, nil
}

// GetCinemaByID retrieves cinema detail with studios
func (u *CinemaUseCase) GetCinemaByID(ctx context.Context, id int) (dto.CinemaResponse, error) {
	cinema, err := u.Repo.Cinema.GetCinemaByID(ctx, id)
	if err != nil {
		return dto.CinemaResponse{}, err
	}

	studios, err := u.Repo.Cinema.GetStudiosByCinemaID(ctx, id)
	if err != nil {
		return dto.CinemaResponse{}, err
	}

	var studioResponses []dto.StudioResponse
	for _, s := range studios {
		studioResponses = append(studioResponses, dto.StudioResponse{
			ID:         s.ID,
			Name:       s.Name,
			TotalSeats: s.TotalSeats,
		})
	}

	return dto.CinemaResponse{
		ID:        cinema.ID,
		Name:      cinema.Name,
		Location:  cinema.Location,
		Studios:   studioResponses,
		CreatedAt: cinema.CreatedAt,
	}, nil
}
