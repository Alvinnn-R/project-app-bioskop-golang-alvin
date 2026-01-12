package usecase

import (
	"context"
	"fmt"
	"session-23/internal/data/entity"
	"session-23/internal/data/repository"

	"session-23/internal/dto"

	"time"
)

type UsecaseServiceCar interface {
	DashboardSerial(ctx context.Context, limit int) (dto.DashboardResponse, error)
	DashboardConcurrent(ctx context.Context, limit int) (dto.DashboardResponse, error)
}

type ServiceCar struct {
	Repo repository.Repository
}

func NewServiceCar(repo *repository.Repository) *ServiceCar {
	return &ServiceCar{Repo: *repo}
}

// DashboardSerial - execute queries sequentially
func (s *ServiceCar) DashboardSerial(ctx context.Context, limit int) (dto.DashboardResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Query 1: Get Latest Cars
	cars, err := s.Repo.RepositoryCar.GetLatestCars(ctx, limit)
	if err != nil {
		return dto.DashboardResponse{}, fmt.Errorf("getLatestCars: %w", err)
	}

	// Query 2: Get Total Cars
	total, err := s.Repo.RepositoryCar.GetTotalCars(ctx)
	if err != nil {
		return dto.DashboardResponse{}, fmt.Errorf("getTotalCars: %w", err)
	}

	// Query 3: Get Price Stats
	stats, err := s.Repo.RepositoryCar.GetPriceStats(ctx)
	if err != nil {
		return dto.DashboardResponse{}, fmt.Errorf("getPriceStats: %w", err)
	}

	return dto.DashboardResponse{
		TotalCars: total,
		Stats:     stats,
		Cars:      cars,
	}, nil
}

// DashboardConcurrent - execute queries concurrently using goroutines
func (s *ServiceCar) DashboardConcurrent(ctx context.Context, limit int) (dto.DashboardResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	carsCh := make(chan dto.ResultCars)
	totalCh := make(chan dto.ResultTotal)
	statsCh := make(chan dto.ResultStats)

	// Run 3 queries concurrently
	go func() {
		cars, err := s.Repo.RepositoryCar.GetLatestCars(ctx, limit)
		carsCh <- dto.ResultCars{Data: cars, Err: err}
	}()
	go func() {
		total, err := s.Repo.RepositoryCar.GetTotalCars(ctx)
		totalCh <- dto.ResultTotal{Data: total, Err: err}
	}()
	go func() {
		stats, err := s.Repo.RepositoryCar.GetPriceStats(ctx)
		statsCh <- dto.ResultStats{Data: stats, Err: err}
	}()

	var (
		cars  []entity.Car
		total int64
		stats dto.PriceStats
	)

	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			return dto.DashboardResponse{}, ctx.Err()

		case r := <-carsCh:
			if r.Err != nil {
				return dto.DashboardResponse{}, fmt.Errorf("getLatestCars: %w", r.Err)
			}
			cars = r.Data

		case r := <-totalCh:
			if r.Err != nil {
				return dto.DashboardResponse{}, fmt.Errorf("getTotalCars: %w", r.Err)
			}
			total = r.Data

		case r := <-statsCh:
			if r.Err != nil {
				return dto.DashboardResponse{}, fmt.Errorf("getPriceStats: %w", r.Err)
			}
			stats = r.Data
		}
	}

	return dto.DashboardResponse{
		TotalCars: total,
		Stats:     stats,
		Cars:      cars,
	}, nil
}
