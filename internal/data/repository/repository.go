package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	RepositoryCar *RepositoryCar
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		RepositoryCar: NewRepositoryCar(db),
	}
}
