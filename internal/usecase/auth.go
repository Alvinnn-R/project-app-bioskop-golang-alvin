package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/pkg/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCaseInterface interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.ResponseUser, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (entity.User, error)
	VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) error
	ResendOTP(ctx context.Context, req dto.ResendOTPRequest) error
}

type AuthUseCase struct {
	Repo         *repository.Repository
	EmailService *utils.EmailService
}

func NewAuthUseCase(repo *repository.Repository) AuthUseCaseInterface {
	return &AuthUseCase{
		Repo:         repo,
		EmailService: utils.NewEmailService(),
	}
}

// generateOTP generates a random 6-digit OTP
func generateOTP() (string, error) {
	const digits = "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}

// Register creates a new user account and sends OTP
func (u *AuthUseCase) Register(ctx context.Context, req dto.RegisterRequest) (dto.ResponseUser, error) {
	// Check if username already exists
	existing, _ := u.Repo.Auth.GetUserByUsername(ctx, req.Username)
	if existing.ID != 0 {
		return dto.ResponseUser{}, errors.New("username already exists")
	}

	// Check if email already exists
	existingEmail, _ := u.Repo.Auth.GetUserByEmail(ctx, req.Email)
	if existingEmail.ID != 0 {
		return dto.ResponseUser{}, errors.New("email already exists")
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

	// Generate and send OTP
	otpCode, err := generateOTP()
	if err != nil {
		return dto.ResponseUser{}, errors.New("failed to generate OTP")
	}

	otp := entity.OTP{
		UserID:    id,
		OTPCode:   otpCode,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	if err := u.Repo.Auth.CreateOTP(ctx, otp); err != nil {
		return dto.ResponseUser{}, errors.New("failed to create OTP")
	}

	// Send OTP via email (goroutine for non-blocking)
	go u.EmailService.SendOTP(req.Email, req.Username, otpCode)

	return dto.ResponseUser{
		ID:         id,
		Username:   req.Username,
		Email:      req.Email,
		IsVerified: false,
	}, nil
}

// VerifyOTP verifies the OTP and marks user as verified
func (u *AuthUseCase) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) error {
	// Get user by email
	user, err := u.Repo.Auth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if already verified
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// Get valid OTP
	otp, err := u.Repo.Auth.GetValidOTP(ctx, user.ID, req.OTP)
	if err != nil {
		return errors.New("invalid or expired OTP")
	}

	// Mark OTP as used
	if err := u.Repo.Auth.MarkOTPUsed(ctx, otp.ID); err != nil {
		return errors.New("failed to verify OTP")
	}

	// Mark user as verified
	if err := u.Repo.Auth.UpdateUserVerified(ctx, user.ID); err != nil {
		return errors.New("failed to verify user")
	}

	return nil
}

// ResendOTP generates and sends a new OTP
func (u *AuthUseCase) ResendOTP(ctx context.Context, req dto.ResendOTPRequest) error {
	// Get user by email
	user, err := u.Repo.Auth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if already verified
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// Invalidate existing OTPs
	if err := u.Repo.Auth.InvalidateUserOTPs(ctx, user.ID); err != nil {
		return errors.New("failed to invalidate existing OTPs")
	}

	// Generate new OTP
	otpCode, err := generateOTP()
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	otp := entity.OTP{
		UserID:    user.ID,
		OTPCode:   otpCode,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	if err := u.Repo.Auth.CreateOTP(ctx, otp); err != nil {
		return errors.New("failed to create OTP")
	}

	// Send OTP via email (goroutine for non-blocking)
	go u.EmailService.SendOTP(req.Email, user.Username, otpCode)

	return nil
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

	// Check if user is verified
	if !user.IsVerified {
		return dto.LoginResponse{}, errors.New("email not verified, please verify your email first")
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
