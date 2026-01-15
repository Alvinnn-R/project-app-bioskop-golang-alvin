package usecase

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// =====================
// Mock Repository untuk Auth
// =====================

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepo) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepo) UpdateUserVerified(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepo) CreateSession(ctx context.Context, session entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockAuthRepo) GetSessionByToken(ctx context.Context, token string) (entity.Session, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(entity.Session), args.Error(1)
}

func (m *MockAuthRepo) RevokeSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthRepo) CreateOTP(ctx context.Context, otp entity.OTP) error {
	args := m.Called(ctx, otp)
	return args.Error(0)
}

func (m *MockAuthRepo) GetValidOTP(ctx context.Context, userID int, otpCode string) (entity.OTP, error) {
	args := m.Called(ctx, userID, otpCode)
	return args.Get(0).(entity.OTP), args.Error(1)
}

func (m *MockAuthRepo) MarkOTPUsed(ctx context.Context, otpID int) error {
	args := m.Called(ctx, otpID)
	return args.Error(0)
}

func (m *MockAuthRepo) InvalidateUserOTPs(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// =====================
// Auth UseCase Tests
// =====================

func TestAuthUseCase_Login_Success(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	now := time.Now()
	// Generate real bcrypt hash for "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := entity.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		IsVerified:   true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	mockAuthRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil)
	mockAuthRepo.On("CreateSession", mock.Anything, mock.AnythingOfType("entity.Session")).Return(nil)

	req := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	result, err := usecase.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "testuser", result.User.Username)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_Login_InvalidUsername(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	mockAuthRepo.On("GetUserByUsername", mock.Anything, "nonexistent").Return(entity.User{}, errors.New("user not found"))

	req := dto.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	_, err := usecase.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid username or password")
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_Login_InvalidPassword(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := entity.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		IsVerified:   true,
		CreatedAt:    now,
	}

	mockAuthRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil)

	req := dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err := usecase.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid username or password")
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_Login_EmailNotVerified(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := entity.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		IsVerified:   false, // Email not verified
		CreatedAt:    now,
	}

	mockAuthRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil)

	req := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	_, err := usecase.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email not verified")
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_Logout_Success(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	mockAuthRepo.On("RevokeSession", mock.Anything, "test-token").Return(nil)

	err := usecase.Logout(context.Background(), "test-token")

	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_Logout_Error(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	mockAuthRepo.On("RevokeSession", mock.Anything, "invalid-token").Return(errors.New("session not found"))

	err := usecase.Logout(context.Background(), "invalid-token")

	assert.Error(t, err)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_ValidateToken_Success(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	now := time.Now()
	session := entity.Session{
		ID:        1,
		UserID:    1,
		Token:     "valid-token",
		ExpiredAt: now.Add(24 * time.Hour),
	}

	user := entity.User{
		ID:         1,
		Username:   "testuser",
		Email:      "test@example.com",
		IsVerified: true,
		CreatedAt:  now,
	}

	mockAuthRepo.On("GetSessionByToken", mock.Anything, "valid-token").Return(session, nil)
	mockAuthRepo.On("GetUserByID", mock.Anything, 1).Return(user, nil)

	result, err := usecase.ValidateToken(context.Background(), "valid-token")

	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "testuser", result.Username)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_ValidateToken_InvalidToken(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	mockAuthRepo.On("GetSessionByToken", mock.Anything, "invalid-token").Return(entity.Session{}, errors.New("token not found"))

	_, err := usecase.ValidateToken(context.Background(), "invalid-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired token")
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_VerifyOTP_Success(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	now := time.Now()
	user := entity.User{
		ID:         1,
		Username:   "testuser",
		Email:      "test@example.com",
		IsVerified: false,
		CreatedAt:  now,
	}

	otp := entity.OTP{
		ID:        1,
		UserID:    1,
		OTPCode:   "123456",
		ExpiredAt: now.Add(5 * time.Minute),
		IsUsed:    false,
	}

	mockAuthRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockAuthRepo.On("GetValidOTP", mock.Anything, 1, "123456").Return(otp, nil)
	mockAuthRepo.On("MarkOTPUsed", mock.Anything, 1).Return(nil)
	mockAuthRepo.On("UpdateUserVerified", mock.Anything, 1).Return(nil)

	req := dto.VerifyOTPRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	err := usecase.VerifyOTP(context.Background(), req)

	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_VerifyOTP_UserNotFound(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	mockAuthRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(entity.User{}, errors.New("not found"))

	req := dto.VerifyOTPRequest{
		Email: "notfound@example.com",
		OTP:   "123456",
	}

	err := usecase.VerifyOTP(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_VerifyOTP_AlreadyVerified(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	now := time.Now()
	user := entity.User{
		ID:         1,
		Username:   "testuser",
		Email:      "test@example.com",
		IsVerified: true, // Already verified
		CreatedAt:  now,
	}

	mockAuthRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

	req := dto.VerifyOTPRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	err := usecase.VerifyOTP(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already verified")
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUseCase_VerifyOTP_InvalidOTP(t *testing.T) {
	mockAuthRepo := new(MockAuthRepo)
	repo := &repository.Repository{Auth: mockAuthRepo}
	usecase := &AuthUseCase{Repo: repo}

	now := time.Now()
	user := entity.User{
		ID:         1,
		Username:   "testuser",
		Email:      "test@example.com",
		IsVerified: false,
		CreatedAt:  now,
	}

	mockAuthRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockAuthRepo.On("GetValidOTP", mock.Anything, 1, "000000").Return(entity.OTP{}, errors.New("invalid otp"))

	req := dto.VerifyOTPRequest{
		Email: "test@example.com",
		OTP:   "000000",
	}

	err := usecase.VerifyOTP(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired OTP")
	mockAuthRepo.AssertExpectations(t)
}
