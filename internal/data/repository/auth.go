package repository

import (
	"context"
	"project-app-bioskop/internal/data/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepoInterface interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUserVerified(ctx context.Context, userID int) error
	CreateSession(ctx context.Context, session entity.Session) error
	GetSessionByToken(ctx context.Context, token string) (entity.Session, error)
	RevokeSession(ctx context.Context, token string) error
	// OTP functions
	CreateOTP(ctx context.Context, otp entity.OTP) error
	GetValidOTP(ctx context.Context, userID int, otpCode string) (entity.OTP, error)
	MarkOTPUsed(ctx context.Context, otpID int) error
	InvalidateUserOTPs(ctx context.Context, userID int) error
}

type AuthRepo struct {
	DB *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) AuthRepoInterface {
	return &AuthRepo{DB: db}
}

// CreateUser inserts a new user into database
func (r *AuthRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := r.DB.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByUsername retrieves user by username
func (r *AuthRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	query := `SELECT id, username, email, password_hash, is_verified, created_at, updated_at 
			  FROM users WHERE username = $1`
	var user entity.User
	err := r.DB.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByID retrieves user by ID
func (r *AuthRepo) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	query := `SELECT id, username, email, password_hash, is_verified, created_at, updated_at 
			  FROM users WHERE id = $1`
	var user entity.User
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// CreateSession creates a new login session
func (r *AuthRepo) CreateSession(ctx context.Context, session entity.Session) error {
	query := `INSERT INTO sessions (user_id, token, expired_at) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(ctx, query, session.UserID, session.Token, session.ExpiredAt)
	return err
}

// GetSessionByToken retrieves active session by token
func (r *AuthRepo) GetSessionByToken(ctx context.Context, token string) (entity.Session, error) {
	query := `SELECT id, user_id, token, expired_at, revoked_at, created_at 
			  FROM sessions WHERE token = $1 AND revoked_at IS NULL AND expired_at > NOW()`
	var session entity.Session
	err := r.DB.QueryRow(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.ExpiredAt, &session.RevokedAt, &session.CreatedAt,
	)
	if err != nil {
		return session, err
	}
	return session, nil
}

// RevokeSession invalidates a session token
func (r *AuthRepo) RevokeSession(ctx context.Context, token string) error {
	query := `UPDATE sessions SET revoked_at = NOW() WHERE token = $1`
	_, err := r.DB.Exec(ctx, query, token)
	return err
}

// GetUserByEmail retrieves user by email
func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `SELECT id, username, email, password_hash, is_verified, created_at, updated_at 
			  FROM users WHERE email = $1`
	var user entity.User
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUserVerified marks user as verified
func (r *AuthRepo) UpdateUserVerified(ctx context.Context, userID int) error {
	query := `UPDATE users SET is_verified = true, updated_at = NOW() WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, userID)
	return err
}

// CreateOTP creates a new OTP record
func (r *AuthRepo) CreateOTP(ctx context.Context, otp entity.OTP) error {
	query := `INSERT INTO otps (user_id, otp_code, expired_at) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(ctx, query, otp.UserID, otp.OTPCode, otp.ExpiredAt)
	return err
}

// GetValidOTP retrieves a valid (not expired, not used) OTP
func (r *AuthRepo) GetValidOTP(ctx context.Context, userID int, otpCode string) (entity.OTP, error) {
	query := `SELECT id, user_id, otp_code, expired_at, is_used, created_at 
			  FROM otps WHERE user_id = $1 AND otp_code = $2 AND is_used = false AND expired_at > NOW()`
	var otp entity.OTP
	err := r.DB.QueryRow(ctx, query, userID, otpCode).Scan(
		&otp.ID, &otp.UserID, &otp.OTPCode, &otp.ExpiredAt, &otp.IsUsed, &otp.CreatedAt,
	)
	if err != nil {
		return otp, err
	}
	return otp, nil
}

// MarkOTPUsed marks an OTP as used
func (r *AuthRepo) MarkOTPUsed(ctx context.Context, otpID int) error {
	query := `UPDATE otps SET is_used = true WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, otpID)
	return err
}

// InvalidateUserOTPs marks all existing OTPs for a user as used
func (r *AuthRepo) InvalidateUserOTPs(ctx context.Context, userID int) error {
	query := `UPDATE otps SET is_used = true WHERE user_id = $1 AND is_used = false`
	_, err := r.DB.Exec(ctx, query, userID)
	return err
}
