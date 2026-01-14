package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository aggregates all repository interfaces
type Repository struct {
	Auth    AuthRepoInterface
	Cinema  CinemaRepoInterface
	Seat    SeatRepoInterface
	Booking BookingRepoInterface
	Payment PaymentRepoInterface
	Movie   MovieRepoInterface
}

// NewRepository creates a new Repository instance with all sub-repositories
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Auth:    NewAuthRepo(db),
		Cinema:  NewCinemaRepo(db),
		Seat:    NewSeatRepo(db),
		Booking: NewBookingRepo(db),
		Payment: NewPaymentRepo(db),
		Movie:   NewMovieRepo(db),
	}
}
