package usecase

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCaseInterface interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.ResponseUser, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (entity.User, error)
}

type AuthUseCase struct {
	Repo *repository.Repository
}

func NewAuthUseCase(repo *repository.Repository) AuthUseCaseInterface {
	return &AuthUseCase{Repo: repo}
}

// Register creates a new user account
func (u *AuthUseCase) Register(ctx context.Context, req dto.RegisterRequest) (dto.ResponseUser, error) {
	// Check if username already exists
	existing, _ := u.Repo.Auth.GetUserByUsername(ctx, req.Username)
	if existing.ID != 0 {
		return dto.ResponseUser{}, errors.New("username already exists")
	}

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.ResponseUser{}, err
	}

	user := entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	id, err := u.Repo.Auth.CreateUser(ctx, user)
	if err != nil {
		return dto.ResponseUser{}, err
	}

	return dto.ResponseUser{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

// Login authenticates user and returns token
func (u *AuthUseCase) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := u.Repo.Auth.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, errors.New("invalid username or password")
	}

	// Generate token
	token := uuid.New().String()
	expiredAt := time.Now().Add(24 * time.Hour)

	session := entity.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: expiredAt,
	}

	if err := u.Repo.Auth.CreateSession(ctx, session); err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token:     token,
		ExpiredAt: expiredAt,
		User: dto.ResponseUser{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt,
		},
	}, nil
}

// Logout invalidates user session
func (u *AuthUseCase) Logout(ctx context.Context, token string) error {
	return u.Repo.Auth.RevokeSession(ctx, token)
}

// ValidateToken validates session token and returns user
func (u *AuthUseCase) ValidateToken(ctx context.Context, token string) (entity.User, error) {
	session, err := u.Repo.Auth.GetSessionByToken(ctx, token)
	if err != nil {
		return entity.User{}, errors.New("invalid or expired token")
	}

	user, err := u.Repo.Auth.GetUserByID(ctx, session.UserID)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
