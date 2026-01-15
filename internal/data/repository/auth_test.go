package repository

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepo_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	user := entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	t.Run("Success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Username, user.Email, user.PasswordHash).
			WillReturnRows(rows)

		id, err := repo.CreateUser(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Username, user.Email, user.PasswordHash).
			WillReturnError(errors.New("database error"))

		id, err := repo.CreateUser(context.Background(), user)
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Contains(t, err.Error(), "database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_GetUserByUsername(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	t.Run("Success - User Found", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows([]string{
			"id", "username", "email", "password_hash", "is_verified", "created_at", "updated_at",
		}).AddRow(1, "testuser", "test@example.com", "hashedpwd", true, now, now)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE username").
			WithArgs("testuser").
			WillReturnRows(rows)

		user, err := repo.GetUserByUsername(context.Background(), "testuser")
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE username").
			WithArgs("nonexistent").
			WillReturnError(pgx.ErrNoRows)

		user, err := repo.GetUserByUsername(context.Background(), "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, 0, user.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_GetUserByEmail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	t.Run("Success - User Found", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows([]string{
			"id", "username", "email", "password_hash", "is_verified", "created_at", "updated_at",
		}).AddRow(1, "testuser", "test@example.com", "hashedpwd", true, now, now)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE email").
			WithArgs("test@example.com").
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(context.Background(), "test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE email").
			WithArgs("nonexistent@example.com").
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetUserByEmail(context.Background(), "nonexistent@example.com")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_GetUserByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	t.Run("Success - User Found", func(t *testing.T) {
		now := time.Now()
		rows := pgxmock.NewRows([]string{
			"id", "username", "email", "password_hash", "is_verified", "created_at", "updated_at",
		}).AddRow(1, "testuser", "test@example.com", "hashedpwd", true, now, now)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE id").
			WithArgs(1).
			WillReturnRows(rows)

		user, err := repo.GetUserByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id").
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetUserByID(context.Background(), 999)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_CreateSession(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	session := entity.Session{
		UserID:    1,
		Token:     "test-token-123",
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sessions").
			WithArgs(session.UserID, session.Token, session.ExpiredAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.CreateSession(context.Background(), session)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sessions").
			WithArgs(session.UserID, session.Token, session.ExpiredAt).
			WillReturnError(errors.New("database error"))

		err := repo.CreateSession(context.Background(), session)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_RevokeSession(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE sessions SET revoked_at").
			WithArgs("test-token").
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.RevokeSession(context.Background(), "test-token")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_UpdateUserVerified(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET is_verified").
			WithArgs(1).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.UpdateUserVerified(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_CreateOTP(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	otp := entity.OTP{
		UserID:    1,
		OTPCode:   "123456",
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO otps").
			WithArgs(otp.UserID, otp.OTPCode, otp.ExpiredAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.CreateOTP(context.Background(), otp)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO otps").
			WithArgs(otp.UserID, otp.OTPCode, otp.ExpiredAt).
			WillReturnError(errors.New("database error"))

		err := repo.CreateOTP(context.Background(), otp)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_MarkOTPUsed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE otps SET is_used").
			WithArgs(1).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.MarkOTPUsed(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAuthRepo_InvalidateUserOTPs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewAuthRepo(mock)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE otps SET is_used").
			WithArgs(1).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.InvalidateUserOTPs(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
