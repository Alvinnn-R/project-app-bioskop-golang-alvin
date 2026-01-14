package repository

import (
	"context"
	"project-app-bioskop/internal/data/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepoInterface interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	CreateSession(ctx context.Context, session entity.Session) error
	GetSessionByToken(ctx context.Context, token string) (entity.Session, error)
	RevokeSession(ctx context.Context, token string) error
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
