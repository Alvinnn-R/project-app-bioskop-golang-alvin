package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestBookingRepo_GetBookingByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewBookingRepo(mock)

	t.Run("Success - Booking Found", func(t *testing.T) {
		now := time.Now()

		rows := pgxmock.NewRows([]string{
			"id", "user_id", "showtime_id", "status", "total_amount", "created_at",
		}).AddRow(1, 1, 1, "confirmed", 100000.0, now)

		mock.ExpectQuery("SELECT (.+) FROM bookings WHERE id").
			WithArgs(1).
			WillReturnRows(rows)

		booking, err := repo.GetBookingByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, booking.ID)
		assert.Equal(t, 1, booking.UserID)
		assert.Equal(t, 100000.0, booking.TotalAmount)
		assert.Equal(t, "confirmed", booking.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM bookings WHERE id").
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		booking, err := repo.GetBookingByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, 0, booking.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM bookings WHERE id").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		booking, err := repo.GetBookingByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, 0, booking.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBookingRepo_GetBookingsByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewBookingRepo(mock)

	t.Run("Success - Multiple Bookings", func(t *testing.T) {
		now := time.Now()

		rows := pgxmock.NewRows([]string{
			"id", "user_id", "showtime_id", "status", "total_amount", "created_at",
		}).
			AddRow(1, 1, 1, "confirmed", 100000.0, now).
			AddRow(2, 1, 2, "pending", 150000.0, now)

		mock.ExpectQuery("SELECT (.+) FROM bookings WHERE user_id").
			WithArgs(1).
			WillReturnRows(rows)

		bookings, err := repo.GetBookingsByUserID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, bookings, 2)
		assert.Equal(t, 100000.0, bookings[0].TotalAmount)
		assert.Equal(t, 150000.0, bookings[1].TotalAmount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty Result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "user_id", "showtime_id", "status", "total_amount", "created_at",
		})

		mock.ExpectQuery("SELECT (.+) FROM bookings WHERE user_id").
			WithArgs(999).
			WillReturnRows(rows)

		bookings, err := repo.GetBookingsByUserID(context.Background(), 999)
		assert.NoError(t, err)
		assert.Len(t, bookings, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM bookings WHERE user_id").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		bookings, err := repo.GetBookingsByUserID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, bookings)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBookingRepo_GetBookingSeats(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewBookingRepo(mock)

	t.Run("Success - Get Booking Seats", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "booking_id", "seat_id", "price_snapshot",
		}).
			AddRow(1, 1, 10, 50000.0).
			AddRow(2, 1, 11, 50000.0)

		mock.ExpectQuery("SELECT (.+) FROM booking_seats WHERE booking_id").
			WithArgs(1).
			WillReturnRows(rows)

		seats, err := repo.GetBookingSeats(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, seats, 2)
		assert.Equal(t, 50000.0, seats[0].PriceSnapshot)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty Result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "booking_id", "seat_id", "price_snapshot",
		})

		mock.ExpectQuery("SELECT (.+) FROM booking_seats WHERE booking_id").
			WithArgs(999).
			WillReturnRows(rows)

		seats, err := repo.GetBookingSeats(context.Background(), 999)
		assert.NoError(t, err)
		assert.Len(t, seats, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM booking_seats WHERE booking_id").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		seats, err := repo.GetBookingSeats(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, seats)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBookingRepo_UpdateBookingStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewBookingRepo(mock)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE bookings SET status").
			WithArgs("confirmed", 1).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.UpdateBookingStatus(context.Background(), 1, "confirmed")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectExec("UPDATE bookings SET status").
			WithArgs("confirmed", 1).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateBookingStatus(context.Background(), 1, "confirmed")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
