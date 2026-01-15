package usecase

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =====================
// Mock Repository untuk Movie
// =====================

type MockMovieRepo struct {
	mock.Mock
}

func (m *MockMovieRepo) GetAllMovies(ctx context.Context, page, limit int) ([]entity.Movie, dto.Pagination, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]entity.Movie), args.Get(1).(dto.Pagination), args.Error(2)
}

func (m *MockMovieRepo) GetMovieByID(ctx context.Context, id int) (entity.Movie, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Movie), args.Error(1)
}

// =====================
// Movie UseCase Tests
// =====================

func TestMovieUseCase_GetAllMovies_Success(t *testing.T) {
	mockMovieRepo := new(MockMovieRepo)
	repo := &repository.Repository{Movie: mockMovieRepo}
	usecase := &MovieUseCase{Repo: repo}

	now := time.Now()
	movies := []entity.Movie{
		{ID: 1, Title: "Avengers", PosterURL: "http://url1", Genres: []string{"Action"}, Rating: 8.5, DurationMinutes: 180, CreatedAt: now},
		{ID: 2, Title: "Batman", PosterURL: "http://url2", Genres: []string{"Action"}, Rating: 8.0, DurationMinutes: 150, CreatedAt: now},
	}
	pagination := dto.Pagination{CurrentPage: 1, TotalPages: 1, TotalRecords: 2, Limit: 10}

	mockMovieRepo.On("GetAllMovies", mock.Anything, 1, 10).Return(movies, pagination, nil)

	result, pag, err := usecase.GetAllMovies(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Avengers", result[0].Title)
	assert.Equal(t, "Batman", result[1].Title)
	assert.Equal(t, 2, pag.TotalRecords)
	mockMovieRepo.AssertExpectations(t)
}

func TestMovieUseCase_GetAllMovies_Error(t *testing.T) {
	mockMovieRepo := new(MockMovieRepo)
	repo := &repository.Repository{Movie: mockMovieRepo}
	usecase := &MovieUseCase{Repo: repo}

	mockMovieRepo.On("GetAllMovies", mock.Anything, 1, 10).Return([]entity.Movie{}, dto.Pagination{}, errors.New("database error"))

	result, _, err := usecase.GetAllMovies(context.Background(), 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockMovieRepo.AssertExpectations(t)
}

func TestMovieUseCase_GetAllMovies_Empty(t *testing.T) {
	mockMovieRepo := new(MockMovieRepo)
	repo := &repository.Repository{Movie: mockMovieRepo}
	usecase := &MovieUseCase{Repo: repo}

	mockMovieRepo.On("GetAllMovies", mock.Anything, 1, 10).Return([]entity.Movie{}, dto.Pagination{TotalRecords: 0}, nil)

	result, pag, err := usecase.GetAllMovies(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, pag.TotalRecords)
	mockMovieRepo.AssertExpectations(t)
}

func TestMovieUseCase_GetMovieByID_Success(t *testing.T) {
	mockMovieRepo := new(MockMovieRepo)
	repo := &repository.Repository{Movie: mockMovieRepo}
	usecase := &MovieUseCase{Repo: repo}

	now := time.Now()
	movie := entity.Movie{
		ID:              1,
		Title:           "Avengers",
		PosterURL:       "http://poster.url",
		Genres:          []string{"Action", "Sci-Fi"},
		Rating:          8.5,
		DurationMinutes: 180,
		CreatedAt:       now,
	}

	mockMovieRepo.On("GetMovieByID", mock.Anything, 1).Return(movie, nil)

	result, err := usecase.GetMovieByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Avengers", result.Title)
	assert.Equal(t, 8.5, result.Rating)
	mockMovieRepo.AssertExpectations(t)
}

func TestMovieUseCase_GetMovieByID_NotFound(t *testing.T) {
	mockMovieRepo := new(MockMovieRepo)
	repo := &repository.Repository{Movie: mockMovieRepo}
	usecase := &MovieUseCase{Repo: repo}

	mockMovieRepo.On("GetMovieByID", mock.Anything, 999).Return(entity.Movie{}, errors.New("movie not found"))

	result, err := usecase.GetMovieByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, 0, result.ID)
	assert.Contains(t, err.Error(), "not found")
	mockMovieRepo.AssertExpectations(t)
}
