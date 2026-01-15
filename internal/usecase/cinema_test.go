package usecase

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =====================
// Mock Repository untuk Cinema
// =====================

type MockCinemaRepo struct {
	mock.Mock
}

func (m *MockCinemaRepo) GetAllCinemas(ctx context.Context, limit, offset int) ([]entity.Cinema, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]entity.Cinema), args.Error(1)
}

func (m *MockCinemaRepo) GetCinemaByID(ctx context.Context, id int) (entity.Cinema, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Cinema), args.Error(1)
}

func (m *MockCinemaRepo) GetStudiosByCinemaID(ctx context.Context, cinemaID int) ([]entity.Studio, error) {
	args := m.Called(ctx, cinemaID)
	return args.Get(0).([]entity.Studio), args.Error(1)
}

func (m *MockCinemaRepo) CountCinemas(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// =====================
// Cinema UseCase Tests
// =====================

func TestCinemaUseCase_GetAllCinemas_Success(t *testing.T) {
	mockCinemaRepo := new(MockCinemaRepo)
	repo := &repository.Repository{Cinema: mockCinemaRepo}
	usecase := &CinemaUseCase{Repo: repo}

	now := time.Now()
	cinemas := []entity.Cinema{
		{ID: 1, Name: "Cinema XXI Jakarta", Location: "Jakarta", CreatedAt: now},
		{ID: 2, Name: "Cinema XXI Bandung", Location: "Bandung", CreatedAt: now},
	}

	mockCinemaRepo.On("GetAllCinemas", mock.Anything, 10, 0).Return(cinemas, nil)
	mockCinemaRepo.On("CountCinemas", mock.Anything).Return(2, nil)

	result, pagination, err := usecase.GetAllCinemas(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Cinema XXI Jakarta", result[0].Name)
	assert.Equal(t, 2, pagination.TotalRecords)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCinemaUseCase_GetAllCinemas_Error(t *testing.T) {
	mockCinemaRepo := new(MockCinemaRepo)
	repo := &repository.Repository{Cinema: mockCinemaRepo}
	usecase := &CinemaUseCase{Repo: repo}

	mockCinemaRepo.On("GetAllCinemas", mock.Anything, 10, 0).Return([]entity.Cinema{}, errors.New("database error"))

	result, _, err := usecase.GetAllCinemas(context.Background(), 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCinemaUseCase_GetAllCinemas_CountError(t *testing.T) {
	mockCinemaRepo := new(MockCinemaRepo)
	repo := &repository.Repository{Cinema: mockCinemaRepo}
	usecase := &CinemaUseCase{Repo: repo}

	now := time.Now()
	cinemas := []entity.Cinema{
		{ID: 1, Name: "Cinema XXI Jakarta", Location: "Jakarta", CreatedAt: now},
	}

	mockCinemaRepo.On("GetAllCinemas", mock.Anything, 10, 0).Return(cinemas, nil)
	mockCinemaRepo.On("CountCinemas", mock.Anything).Return(0, errors.New("count error"))

	result, _, err := usecase.GetAllCinemas(context.Background(), 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCinemaUseCase_GetAllCinemas_InvalidPage(t *testing.T) {
	mockCinemaRepo := new(MockCinemaRepo)
	repo := &repository.Repository{Cinema: mockCinemaRepo}
	usecase := &CinemaUseCase{Repo: repo}

	now := time.Now()
	cinemas := []entity.Cinema{
		{ID: 1, Name: "Cinema XXI Jakarta", Location: "Jakarta", CreatedAt: now},
	}

	// When page < 1, it should default to page 1 (offset 0)
	mockCinemaRepo.On("GetAllCinemas", mock.Anything, 10, 0).Return(cinemas, nil)
	mockCinemaRepo.On("CountCinemas", mock.Anything).Return(1, nil)

	result, pagination, err := usecase.GetAllCinemas(context.Background(), -1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, pagination.CurrentPage)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCinemaUseCase_GetCinemaByID_Success(t *testing.T) {
	mockCinemaRepo := new(MockCinemaRepo)
	repo := &repository.Repository{Cinema: mockCinemaRepo}
	usecase := &CinemaUseCase{Repo: repo}

	now := time.Now()
	cinema := entity.Cinema{
		ID:        1,
		Name:      "Cinema XXI Jakarta",
		Location:  "Jakarta Selatan",
		CreatedAt: now,
	}

	studios := []entity.Studio{
		{ID: 1, CinemaID: 1, Name: "Studio 1", TotalSeats: 100},
		{ID: 2, CinemaID: 1, Name: "Studio 2", TotalSeats: 150},
	}

	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 1).Return(cinema, nil)
	mockCinemaRepo.On("GetStudiosByCinemaID", mock.Anything, 1).Return(studios, nil)

	result, err := usecase.GetCinemaByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Cinema XXI Jakarta", result.Name)
	assert.Len(t, result.Studios, 2)
	assert.Equal(t, "Studio 1", result.Studios[0].Name)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCinemaUseCase_GetCinemaByID_NotFound(t *testing.T) {
	mockCinemaRepo := new(MockCinemaRepo)
	repo := &repository.Repository{Cinema: mockCinemaRepo}
	usecase := &CinemaUseCase{Repo: repo}

	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 999).Return(entity.Cinema{}, errors.New("cinema not found"))

	result, err := usecase.GetCinemaByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, 0, result.ID)
	mockCinemaRepo.AssertExpectations(t)
}

func TestCinemaUseCase_GetCinemaByID_StudiosError(t *testing.T) {
	mockCinemaRepo := new(MockCinemaRepo)
	repo := &repository.Repository{Cinema: mockCinemaRepo}
	usecase := &CinemaUseCase{Repo: repo}

	now := time.Now()
	cinema := entity.Cinema{
		ID:        1,
		Name:      "Cinema XXI Jakarta",
		Location:  "Jakarta Selatan",
		CreatedAt: now,
	}

	mockCinemaRepo.On("GetCinemaByID", mock.Anything, 1).Return(cinema, nil)
	mockCinemaRepo.On("GetStudiosByCinemaID", mock.Anything, 1).Return([]entity.Studio{}, errors.New("studios error"))

	result, err := usecase.GetCinemaByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, 0, result.ID)
	mockCinemaRepo.AssertExpectations(t)
}
