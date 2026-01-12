package repository

import (
	"context"
	"session-23/internal/data/entity"
	"session-23/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryCar struct {
	db *pgxpool.Pool
}

func NewRepositoryCar(db *pgxpool.Pool) *RepositoryCar {
	return &RepositoryCar{db: db}
}

func (s *RepositoryCar) GetLatestCars(ctx context.Context, limit int) ([]entity.Car, error) {
	const q = `
		SELECT id, brand, model, year, price, mileage, color, condition, first_name, last_name, address, country
		FROM cars
		ORDER BY id DESC
		LIMIT $1;
	`

	rows, err := s.db.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]entity.Car, 0, limit)
	for rows.Next() {
		var c entity.Car
		if err := rows.Scan(
			&c.ID, &c.Brand, &c.Model, &c.Year, &c.Price, &c.Mileage, &c.Color, &c.Condition,
			&c.FirstName, &c.LastName, &c.Address, &c.Country,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *RepositoryCar) GetTotalCars(ctx context.Context) (int64, error) {
	const q = `SELECT COUNT(*) FROM cars;`
	var total int64
	if err := s.db.QueryRow(ctx, q).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (s *RepositoryCar) GetPriceStats(ctx context.Context) (dto.PriceStats, error) {
	const q = `
		SELECT COALESCE(MIN(price),0), COALESCE(MAX(price),0), COALESCE(AVG(price),0)
		FROM cars;
	`
	var st dto.PriceStats
	if err := s.db.QueryRow(ctx, q).Scan(&st.Min, &st.Max, &st.Avg); err != nil {
		return dto.PriceStats{}, err
	}
	return st, nil
}
