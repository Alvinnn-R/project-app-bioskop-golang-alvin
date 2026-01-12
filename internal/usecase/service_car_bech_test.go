package usecase

import (
	"context"
	"os"
	"session-23/internal/data/repository"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var svc *ServiceCar

func TestMain(m *testing.M) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic("DB_DSN required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	svc = NewServiceCar(&repo)
	os.Exit(m.Run())
}

func BenchmarkDashboardSerial(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.DashboardSerial(ctx, 20)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDashboardConcurrent(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.DashboardConcurrent(ctx, 20)
		if err != nil {
			b.Fatal(err)
		}
	}
}
