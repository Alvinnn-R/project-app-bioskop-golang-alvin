package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestSeatRepo_GetSeatsByShowtime(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSeatRepo(mock)

	t.Run("Success - Get Seats", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "seat_code", "studio_id", "is_booked",
		}).
			AddRow(1, "A1", 1, false).
			AddRow(2, "A2", 1, true).
			AddRow(3, "A3", 1, false)

		mock.ExpectQuery("SELECT DISTINCT").
			WithArgs(1).
			WillReturnRows(rows)

		seats, err := repo.GetSeatsByShowtime(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, seats, 3)
		assert.Equal(t, "A1", seats[0].SeatCode)
		assert.False(t, seats[0].IsBooked)
		assert.True(t, seats[1].IsBooked)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty Result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "seat_code", "studio_id", "is_booked",
		})

		mock.ExpectQuery("SELECT DISTINCT").
			WithArgs(999).
			WillReturnRows(rows)

		seats, err := repo.GetSeatsByShowtime(context.Background(), 999)
		assert.NoError(t, err)
		assert.Len(t, seats, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		seats, err := repo.GetSeatsByShowtime(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, seats)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSeatRepo_GetShowtimeByParams(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSeatRepo(mock)

	t.Run("Success - Showtime Found", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "cinema_id", "studio_id", "movie_id", "show_date", "show_time", "price",
		}).AddRow(1, 1, 1, 1, "2024-01-15", "19:00", 50000.0)

		mock.ExpectQuery("SELECT (.+) FROM showtimes WHERE cinema_id").
			WithArgs(1, "2024-01-15", "19:00").
			WillReturnRows(rows)

		showtime, err := repo.GetShowtimeByParams(context.Background(), 1, "2024-01-15", "19:00")
		assert.NoError(t, err)
		assert.Equal(t, 1, showtime.ID)
		assert.Equal(t, 50000.0, showtime.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM showtimes WHERE cinema_id").
			WithArgs(999, "2024-01-15", "19:00").
			WillReturnError(pgx.ErrNoRows)

		showtime, err := repo.GetShowtimeByParams(context.Background(), 999, "2024-01-15", "19:00")
		assert.Error(t, err)
		assert.Equal(t, 0, showtime.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSeatRepo_GetShowtimesByCinema(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSeatRepo(mock)

	t.Run("Success - Get Showtimes", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "cinema_id", "studio_id", "movie_id", "show_date", "show_time", "price",
			"m_id", "m_title", "m_poster_url", "m_genres", "m_rating", "m_duration",
			"s_id", "s_name", "s_total_seats",
		}).
			AddRow(1, 1, 1, 1, "2024-01-15", "19:00", 50000.0,
				1, "Movie 1", "http://poster1.jpg", []string{"Action"}, 8.5, 120,
				1, "Studio 1", 100).
			AddRow(2, 1, 2, 2, "2024-01-15", "21:00", 60000.0,
				2, "Movie 2", "http://poster2.jpg", []string{"Comedy"}, 7.0, 90,
				2, "Studio 2", 150)

		mock.ExpectQuery("SELECT (.+) FROM showtimes").
			WithArgs(1).
			WillReturnRows(rows)

		showtimes, err := repo.GetShowtimesByCinema(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, showtimes, 2)
		assert.Equal(t, 50000.0, showtimes[0].Price)
		assert.Equal(t, "Movie 1", showtimes[0].Movie.Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty Result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "cinema_id", "studio_id", "movie_id", "show_date", "show_time", "price",
			"m_id", "m_title", "m_poster_url", "m_genres", "m_rating", "m_duration",
			"s_id", "s_name", "s_total_seats",
		})

		mock.ExpectQuery("SELECT (.+) FROM showtimes").
			WithArgs(999).
			WillReturnRows(rows)

		showtimes, err := repo.GetShowtimesByCinema(context.Background(), 999)
		assert.NoError(t, err)
		assert.Len(t, showtimes, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM showtimes").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		showtimes, err := repo.GetShowtimesByCinema(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, showtimes)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSeatRepo_GetSeatsByIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSeatRepo(mock)

	t.Run("Success - Get Seats by IDs", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "studio_id", "seat_code",
		}).
			AddRow(1, 1, "A1").
			AddRow(2, 1, "A2")

		mock.ExpectQuery("SELECT (.+) FROM seats WHERE id = ANY").
			WithArgs([]int{1, 2}).
			WillReturnRows(rows)

		seats, err := repo.GetSeatsByIDs(context.Background(), []int{1, 2})
		assert.NoError(t, err)
		assert.Len(t, seats, 2)
		assert.Equal(t, "A1", seats[0].SeatCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty IDs", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "studio_id", "seat_code",
		})

		mock.ExpectQuery("SELECT (.+) FROM seats WHERE id = ANY").
			WithArgs([]int{}).
			WillReturnRows(rows)

		seats, err := repo.GetSeatsByIDs(context.Background(), []int{})
		assert.NoError(t, err)
		assert.Len(t, seats, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM seats WHERE id = ANY").
			WithArgs([]int{1, 2}).
			WillReturnError(errors.New("database error"))

		seats, err := repo.GetSeatsByIDs(context.Background(), []int{1, 2})
		assert.Error(t, err)
		assert.Nil(t, seats)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSeatRepo_CheckSeatsAvailable(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSeatRepo(mock)

	t.Run("All Seats Available", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1, []int{1, 2, 3}).
			WillReturnRows(rows)

		available, err := repo.CheckSeatsAvailable(context.Background(), 1, []int{1, 2, 3})
		assert.NoError(t, err)
		assert.True(t, available)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Some Seats Not Available", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"count"}).AddRow(2)

		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1, []int{1, 2, 3}).
			WillReturnRows(rows)

		available, err := repo.CheckSeatsAvailable(context.Background(), 1, []int{1, 2, 3})
		assert.NoError(t, err)
		assert.False(t, available)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1, []int{1, 2}).
			WillReturnError(errors.New("database error"))

		available, err := repo.CheckSeatsAvailable(context.Background(), 1, []int{1, 2})
		assert.Error(t, err)
		assert.False(t, available)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
