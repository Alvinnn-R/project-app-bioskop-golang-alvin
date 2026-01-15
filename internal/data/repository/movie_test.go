package repository

import (
	"context"
	"errors"
	"project-app-bioskop/internal/dto"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestMovieRepo_GetAllMovies(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewMovieRepo(mock)

	t.Run("Success - Get All Movies", func(t *testing.T) {
		now := time.Now()
		releaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		// Mock count query
		countRows := pgxmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

		// Mock movies query
		rows := pgxmock.NewRows([]string{
			"id", "title", "poster_url", "genres", "rating", "review_count",
			"release_date", "duration_in_minutes", "release_status", "created_at", "updated_at",
		}).
			AddRow(1, "Movie 1", "http://poster1.jpg", []string{"Action"}, 8.5, 100, releaseDate, 120, "now_showing", now, now).
			AddRow(2, "Movie 2", "http://poster2.jpg", []string{"Comedy"}, 7.0, 50, releaseDate, 90, "coming_soon", now, now)

		mock.ExpectQuery("SELECT (.+) FROM movies").
			WithArgs(10, 0).
			WillReturnRows(rows)

		movies, pagination, err := repo.GetAllMovies(context.Background(), 1, 10)
		assert.NoError(t, err)
		assert.Len(t, movies, 2)
		assert.Equal(t, "Movie 1", movies[0].Title)
		assert.Equal(t, "Movie 2", movies[1].Title)
		assert.Equal(t, 2, pagination.TotalRecords)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Count Query Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WillReturnError(errors.New("count error"))

		movies, pagination, err := repo.GetAllMovies(context.Background(), 1, 10)
		assert.Error(t, err)
		assert.Nil(t, movies)
		assert.Equal(t, dto.Pagination{}, pagination)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Query Error", func(t *testing.T) {
		countRows := pgxmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

		mock.ExpectQuery("SELECT (.+) FROM movies").
			WithArgs(10, 0).
			WillReturnError(errors.New("query error"))

		movies, pagination, err := repo.GetAllMovies(context.Background(), 1, 10)
		assert.Error(t, err)
		assert.Nil(t, movies)
		assert.Equal(t, dto.Pagination{}, pagination)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMovieRepo_GetMovieByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewMovieRepo(mock)

	t.Run("Success - Movie Found", func(t *testing.T) {
		now := time.Now()
		releaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		rows := pgxmock.NewRows([]string{
			"id", "title", "poster_url", "genres", "rating", "review_count",
			"release_date", "duration_in_minutes", "release_status", "created_at", "updated_at",
		}).AddRow(1, "Test Movie", "http://poster.jpg", []string{"Action"}, 8.5, 100, releaseDate, 120, "now_showing", now, now)

		mock.ExpectQuery("SELECT (.+) FROM movies WHERE id").
			WithArgs(1).
			WillReturnRows(rows)

		movie, err := repo.GetMovieByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, movie.ID)
		assert.Equal(t, "Test Movie", movie.Title)
		assert.Equal(t, 120, movie.DurationMinutes)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM movies WHERE id").
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		movie, err := repo.GetMovieByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, 0, movie.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM movies WHERE id").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		movie, err := repo.GetMovieByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, 0, movie.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
