package usecase

import (
	"context"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
)

type UserUseCaseInterface interface {
	GetUserByID(ctx context.Context, userID int) (entity.User, error)
	GetUserProfile(ctx context.Context, userID int) (entity.User, error)
}

type UserUseCase struct {
	Repo *repository.Repository
}

func NewUserUseCase(repo *repository.Repository) UserUseCaseInterface {
	return &UserUseCase{Repo: repo}
}

// GetUserByID retrieves user by ID
func (u *UserUseCase) GetUserByID(ctx context.Context, userID int) (entity.User, error) {
	user, err := u.Repo.Auth.GetUserByID(ctx, userID)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// GetUserProfile retrieves user profile
func (u *UserUseCase) GetUserProfile(ctx context.Context, userID int) (entity.User, error) {
	return u.GetUserByID(ctx, userID)
}
