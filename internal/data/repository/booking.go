package repository

import (
	"context"
	"project-app-bioskop/internal/data/entity"

	"github.com/jackc/pgx/v5"
)

type BookingRepoInterface interface {
	CreateBooking(ctx context.Context, booking entity.Booking, seatIDs []int, price float64) (int, error)
	GetBookingByID(ctx context.Context, id int) (entity.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID int) ([]entity.Booking, error)
	GetBookingSeats(ctx context.Context, bookingID int) ([]entity.BookingSeat, error)
	UpdateBookingStatus(ctx context.Context, bookingID int, status string) error
}

type BookingRepo struct {
	DB DBPool
}

func NewBookingRepo(db DBPool) BookingRepoInterface {
	return &BookingRepo{DB: db}
}

// CreateBooking creates a booking with seats in a transaction
func (r *BookingRepo) CreateBooking(ctx context.Context, booking entity.Booking, seatIDs []int, price float64) (int, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Insert booking
	var bookingID int
	totalAmount := price * float64(len(seatIDs))
	query := `INSERT INTO bookings (user_id, showtime_id, status, total_amount) 
			  VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(ctx, query, booking.UserID, booking.ShowtimeID, "pending", totalAmount).Scan(&bookingID)
	if err != nil {
		return 0, err
	}

	// Insert booking seats
	for _, seatID := range seatIDs {
		seatQuery := `INSERT INTO booking_seats (booking_id, seat_id, price_snapshot) VALUES ($1, $2, $3)`
		_, err = tx.Exec(ctx, seatQuery, bookingID, seatID, price)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return bookingID, nil
}

// GetBookingByID retrieves booking by ID
func (r *BookingRepo) GetBookingByID(ctx context.Context, id int) (entity.Booking, error) {
	query := `SELECT id, user_id, showtime_id, status, total_amount, created_at 
			  FROM bookings WHERE id = $1`
	var b entity.Booking
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.ShowtimeID, &b.Status, &b.TotalAmount, &b.CreatedAt,
	)
	if err != nil {
		return b, err
	}
	return b, nil
}

// GetBookingsByUserID retrieves all bookings for a user
func (r *BookingRepo) GetBookingsByUserID(ctx context.Context, userID int) ([]entity.Booking, error) {
	query := `SELECT id, user_id, showtime_id, status, total_amount, created_at 
			  FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []entity.Booking
	for rows.Next() {
		var b entity.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.ShowtimeID, &b.Status, &b.TotalAmount, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

// GetBookingSeats retrieves all seats for a booking
func (r *BookingRepo) GetBookingSeats(ctx context.Context, bookingID int) ([]entity.BookingSeat, error) {
	query := `SELECT id, booking_id, seat_id, price_snapshot FROM booking_seats WHERE booking_id = $1`
	rows, err := r.DB.Query(ctx, query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []entity.BookingSeat
	for rows.Next() {
		var s entity.BookingSeat
		if err := rows.Scan(&s.ID, &s.BookingID, &s.SeatID, &s.PriceSnapshot); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

// UpdateBookingStatus updates booking status
func (r *BookingRepo) UpdateBookingStatus(ctx context.Context, bookingID int, status string) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(ctx, query, status, bookingID)
	return err
}
