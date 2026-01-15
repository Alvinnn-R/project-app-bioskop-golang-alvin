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

func TestCinemaRepo_GetAllCinemas(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCinemaRepo(mock)

	t.Run("Success - Get All Cinemas", func(t *testing.T) {
		now := time.Now()

		rows := pgxmock.NewRows([]string{
			"id", "name", "location", "created_at",
		}).
			AddRow(1, "Cinema 1", "Location 1", now).
			AddRow(2, "Cinema 2", "Location 2", now)

		mock.ExpectQuery("SELECT (.+) FROM cinemas").
			WithArgs(10, 0).
			WillReturnRows(rows)

		cinemas, err := repo.GetAllCinemas(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.Len(t, cinemas, 2)
		assert.Equal(t, "Cinema 1", cinemas[0].Name)
		assert.Equal(t, "Cinema 2", cinemas[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Query Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM cinemas").
			WithArgs(10, 0).
			WillReturnError(errors.New("query error"))

		cinemas, err := repo.GetAllCinemas(context.Background(), 10, 0)
		assert.Error(t, err)
		assert.Nil(t, cinemas)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty Result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "name", "location", "created_at",
		})

		mock.ExpectQuery("SELECT (.+) FROM cinemas").
			WithArgs(10, 0).
			WillReturnRows(rows)

		cinemas, err := repo.GetAllCinemas(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.Len(t, cinemas, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCinemaRepo_GetCinemaByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCinemaRepo(mock)

	t.Run("Success - Cinema Found", func(t *testing.T) {
		now := time.Now()

		rows := pgxmock.NewRows([]string{
			"id", "name", "location", "created_at",
		}).AddRow(1, "Test Cinema", "Test Location", now)

		mock.ExpectQuery("SELECT (.+) FROM cinemas WHERE id").
			WithArgs(1).
			WillReturnRows(rows)

		cinema, err := repo.GetCinemaByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, cinema.ID)
		assert.Equal(t, "Test Cinema", cinema.Name)
		assert.Equal(t, "Test Location", cinema.Location)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM cinemas WHERE id").
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		cinema, err := repo.GetCinemaByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, 0, cinema.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM cinemas WHERE id").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		cinema, err := repo.GetCinemaByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, 0, cinema.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCinemaRepo_GetStudiosByCinemaID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCinemaRepo(mock)

	t.Run("Success - Studios Found", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "cinema_id", "name", "total_seats",
		}).
			AddRow(1, 1, "Studio 1", 100).
			AddRow(2, 1, "Studio 2", 150)

		mock.ExpectQuery("SELECT (.+) FROM studios WHERE cinema_id").
			WithArgs(1).
			WillReturnRows(rows)

		studios, err := repo.GetStudiosByCinemaID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, studios, 2)
		assert.Equal(t, "Studio 1", studios[0].Name)
		assert.Equal(t, 100, studios[0].TotalSeats)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty Result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "cinema_id", "name", "total_seats",
		})

		mock.ExpectQuery("SELECT (.+) FROM studios WHERE cinema_id").
			WithArgs(999).
			WillReturnRows(rows)

		studios, err := repo.GetStudiosByCinemaID(context.Background(), 999)
		assert.NoError(t, err)
		assert.Len(t, studios, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM studios WHERE cinema_id").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		studios, err := repo.GetStudiosByCinemaID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, studios)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCinemaRepo_CountCinemas(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCinemaRepo(mock)

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(rows)

		count, err := repo.CountCinemas(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WillReturnError(errors.New("database error"))

		count, err := repo.CountCinemas(context.Background())
		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
