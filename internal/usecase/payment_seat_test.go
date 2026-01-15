package usecase

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =====================
// Mock Repository untuk Payment
// =====================

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) GetAllPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.PaymentMethod), args.Error(1)
}

func (m *MockPaymentRepo) GetPaymentMethodByID(ctx context.Context, id int) (entity.PaymentMethod, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.PaymentMethod), args.Error(1)
}

func (m *MockPaymentRepo) CreatePayment(ctx context.Context, payment entity.Payment) (int, error) {
	args := m.Called(ctx, payment)
	return args.Int(0), args.Error(1)
}

func (m *MockPaymentRepo) GetPaymentByBookingID(ctx context.Context, bookingID int) (entity.Payment, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).(entity.Payment), args.Error(1)
}

func (m *MockPaymentRepo) UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error {
	args := m.Called(ctx, paymentID, status)
	return args.Error(0)
}

// =====================
// Mock Repository untuk Seat
// =====================

type MockSeatRepo struct {
	mock.Mock
}

func (m *MockSeatRepo) GetSeatsByShowtime(ctx context.Context, showtimeID int) ([]entity.SeatAvailability, error) {
	args := m.Called(ctx, showtimeID)
	return args.Get(0).([]entity.SeatAvailability), args.Error(1)
}

func (m *MockSeatRepo) GetShowtimeByParams(ctx context.Context, cinemaID int, date, time string) (entity.Showtime, error) {
	args := m.Called(ctx, cinemaID, date, time)
	return args.Get(0).(entity.Showtime), args.Error(1)
}

func (m *MockSeatRepo) GetShowtimeByID(ctx context.Context, id int) (entity.Showtime, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Showtime), args.Error(1)
}

func (m *MockSeatRepo) GetSeatsByIDs(ctx context.Context, seatIDs []int) ([]entity.Seat, error) {
	args := m.Called(ctx, seatIDs)
	return args.Get(0).([]entity.Seat), args.Error(1)
}

func (m *MockSeatRepo) CheckSeatsAvailable(ctx context.Context, showtimeID int, seatIDs []int) (bool, error) {
	args := m.Called(ctx, showtimeID, seatIDs)
	return args.Bool(0), args.Error(1)
}

func (m *MockSeatRepo) GetShowtimesByCinema(ctx context.Context, cinemaID int) ([]entity.Showtime, error) {
	args := m.Called(ctx, cinemaID)
	return args.Get(0).([]entity.Showtime), args.Error(1)
}

// =====================
// Mock Repository untuk Booking (untuk Payment test)
// =====================

type MockBookingRepoForPayment struct {
	mock.Mock
}

func (m *MockBookingRepoForPayment) CreateBooking(ctx context.Context, booking entity.Booking, seatIDs []int, price float64) (int, error) {
	args := m.Called(ctx, booking, seatIDs, price)
	return args.Int(0), args.Error(1)
}

func (m *MockBookingRepoForPayment) GetBookingByID(ctx context.Context, id int) (entity.Booking, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Booking), args.Error(1)
}

func (m *MockBookingRepoForPayment) GetBookingsByUserID(ctx context.Context, userID int) ([]entity.Booking, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Booking), args.Error(1)
}

func (m *MockBookingRepoForPayment) GetBookingSeats(ctx context.Context, bookingID int) ([]entity.BookingSeat, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).([]entity.BookingSeat), args.Error(1)
}

func (m *MockBookingRepoForPayment) UpdateBookingStatus(ctx context.Context, bookingID int, status string) error {
	args := m.Called(ctx, bookingID, status)
	return args.Error(0)
}

// =====================
// Payment UseCase Tests
// =====================

func TestPaymentUseCase_GetPaymentMethods_Success(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	repo := &repository.Repository{Payment: mockPaymentRepo}
	usecase := &PaymentUseCase{Repo: repo}

	methods := []entity.PaymentMethod{
		{ID: 1, Name: "Credit Card"},
		{ID: 2, Name: "Bank Transfer"},
		{ID: 3, Name: "E-Wallet"},
	}

	mockPaymentRepo.On("GetAllPaymentMethods", mock.Anything).Return(methods, nil)

	result, err := usecase.GetPaymentMethods(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Credit Card", result[0].Name)
	assert.Equal(t, "Bank Transfer", result[1].Name)
	mockPaymentRepo.AssertExpectations(t)
}

func TestPaymentUseCase_GetPaymentMethods_Error(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	repo := &repository.Repository{Payment: mockPaymentRepo}
	usecase := &PaymentUseCase{Repo: repo}

	mockPaymentRepo.On("GetAllPaymentMethods", mock.Anything).Return([]entity.PaymentMethod{}, errors.New("database error"))

	result, err := usecase.GetPaymentMethods(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	mockPaymentRepo.AssertExpectations(t)
}

func TestPaymentUseCase_GetPaymentMethods_Empty(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	repo := &repository.Repository{Payment: mockPaymentRepo}
	usecase := &PaymentUseCase{Repo: repo}

	mockPaymentRepo.On("GetAllPaymentMethods", mock.Anything).Return([]entity.PaymentMethod{}, nil)

	result, err := usecase.GetPaymentMethods(context.Background())

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockPaymentRepo.AssertExpectations(t)
}

func TestPaymentUseCase_ProcessPayment_BookingNotFound(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	mockBookingRepo := new(MockBookingRepoForPayment)
	repo := &repository.Repository{
		Payment: mockPaymentRepo,
		Booking: mockBookingRepo,
	}
	usecase := &PaymentUseCase{Repo: repo}

	mockBookingRepo.On("GetBookingByID", mock.Anything, 999).Return(entity.Booking{}, errors.New("not found"))

	req := dto.PayRequest{
		BookingID:     999,
		PaymentMethod: 1,
	}

	_, err := usecase.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "booking not found")
	mockBookingRepo.AssertExpectations(t)
}

func TestPaymentUseCase_ProcessPayment_Unauthorized(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	mockBookingRepo := new(MockBookingRepoForPayment)
	repo := &repository.Repository{
		Payment: mockPaymentRepo,
		Booking: mockBookingRepo,
	}
	usecase := &PaymentUseCase{Repo: repo}

	booking := entity.Booking{
		ID:     1,
		UserID: 2, // Different user
		Status: "pending",
	}

	mockBookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)

	req := dto.PayRequest{
		BookingID:     1,
		PaymentMethod: 1,
	}

	_, err := usecase.ProcessPayment(context.Background(), 1, req) // UserID 1 trying to pay for UserID 2's booking

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockBookingRepo.AssertExpectations(t)
}

func TestPaymentUseCase_ProcessPayment_AlreadyPaid(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	mockBookingRepo := new(MockBookingRepoForPayment)
	repo := &repository.Repository{
		Payment: mockPaymentRepo,
		Booking: mockBookingRepo,
	}
	usecase := &PaymentUseCase{Repo: repo}

	booking := entity.Booking{
		ID:     1,
		UserID: 1,
		Status: "paid", // Already paid
	}

	mockBookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)

	req := dto.PayRequest{
		BookingID:     1,
		PaymentMethod: 1,
	}

	_, err := usecase.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already paid")
	mockBookingRepo.AssertExpectations(t)
}

func TestPaymentUseCase_ProcessPayment_BookingCancelled(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	mockBookingRepo := new(MockBookingRepoForPayment)
	repo := &repository.Repository{
		Payment: mockPaymentRepo,
		Booking: mockBookingRepo,
	}
	usecase := &PaymentUseCase{Repo: repo}

	booking := entity.Booking{
		ID:     1,
		UserID: 1,
		Status: "cancelled",
	}

	mockBookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)

	req := dto.PayRequest{
		BookingID:     1,
		PaymentMethod: 1,
	}

	_, err := usecase.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cancelled")
	mockBookingRepo.AssertExpectations(t)
}

func TestPaymentUseCase_ProcessPayment_InvalidPaymentMethod(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	mockBookingRepo := new(MockBookingRepoForPayment)
	repo := &repository.Repository{
		Payment: mockPaymentRepo,
		Booking: mockBookingRepo,
	}
	usecase := &PaymentUseCase{Repo: repo}

	booking := entity.Booking{
		ID:     1,
		UserID: 1,
		Status: "pending",
	}

	mockBookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)
	mockPaymentRepo.On("GetPaymentMethodByID", mock.Anything, 999).Return(entity.PaymentMethod{}, errors.New("not found"))

	req := dto.PayRequest{
		BookingID:     1,
		PaymentMethod: 999,
	}

	_, err := usecase.ProcessPayment(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payment method")
	mockBookingRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

// =====================
// Seat UseCase Tests
// =====================

func TestSeatUseCase_GetSeatAvailability_Success(t *testing.T) {
	mockSeatRepo := new(MockSeatRepo)
	repo := &repository.Repository{Seat: mockSeatRepo}
	usecase := &SeatUseCase{Repo: repo}

	showtime := entity.Showtime{
		ID:       1,
		CinemaID: 1,
		ShowDate: "2026-01-15",
		ShowTime: "19:00",
	}

	seats := []entity.SeatAvailability{
		{ID: 1, SeatCode: "A1", StudioID: 1, IsBooked: false},
		{ID: 2, SeatCode: "A2", StudioID: 1, IsBooked: true},
		{ID: 3, SeatCode: "A3", StudioID: 1, IsBooked: false},
	}

	mockSeatRepo.On("GetShowtimeByParams", mock.Anything, 1, "2026-01-15", "19:00").Return(showtime, nil)
	mockSeatRepo.On("GetSeatsByShowtime", mock.Anything, 1).Return(seats, nil)

	result, err := usecase.GetSeatAvailability(context.Background(), 1, "2026-01-15", "19:00")

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "A1", result[0].SeatCode)
	assert.False(t, result[0].IsBooked)
	assert.True(t, result[1].IsBooked)
	mockSeatRepo.AssertExpectations(t)
}

func TestSeatUseCase_GetSeatAvailability_ShowtimeNotFound(t *testing.T) {
	mockSeatRepo := new(MockSeatRepo)
	repo := &repository.Repository{Seat: mockSeatRepo}
	usecase := &SeatUseCase{Repo: repo}

	mockSeatRepo.On("GetShowtimeByParams", mock.Anything, 1, "2026-01-15", "19:00").Return(entity.Showtime{}, errors.New("showtime not found"))

	result, err := usecase.GetSeatAvailability(context.Background(), 1, "2026-01-15", "19:00")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockSeatRepo.AssertExpectations(t)
}

func TestSeatUseCase_GetSeatAvailability_SeatsError(t *testing.T) {
	mockSeatRepo := new(MockSeatRepo)
	repo := &repository.Repository{Seat: mockSeatRepo}
	usecase := &SeatUseCase{Repo: repo}

	showtime := entity.Showtime{
		ID:       1,
		CinemaID: 1,
		ShowDate: "2026-01-15",
		ShowTime: "19:00",
	}

	mockSeatRepo.On("GetShowtimeByParams", mock.Anything, 1, "2026-01-15", "19:00").Return(showtime, nil)
	mockSeatRepo.On("GetSeatsByShowtime", mock.Anything, 1).Return([]entity.SeatAvailability{}, errors.New("database error"))

	result, err := usecase.GetSeatAvailability(context.Background(), 1, "2026-01-15", "19:00")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockSeatRepo.AssertExpectations(t)
}

func TestSeatUseCase_GetShowtimesByCinema_Success(t *testing.T) {
	mockSeatRepo := new(MockSeatRepo)
	repo := &repository.Repository{Seat: mockSeatRepo}
	usecase := &SeatUseCase{Repo: repo}

	movie := &entity.Movie{ID: 1, Title: "Avengers", Genres: []string{"Action"}, Rating: 8.5, DurationMinutes: 180}
	studio := &entity.Studio{ID: 1, Name: "Studio 1", TotalSeats: 100}

	showtimes := []entity.Showtime{
		{ID: 1, CinemaID: 1, ShowDate: "2026-01-15", ShowTime: "14:00", Price: 50000, Movie: movie, Studio: studio},
		{ID: 2, CinemaID: 1, ShowDate: "2026-01-15", ShowTime: "19:00", Price: 60000, Movie: movie, Studio: studio},
	}

	mockSeatRepo.On("GetShowtimesByCinema", mock.Anything, 1).Return(showtimes, nil)

	result, err := usecase.GetShowtimesByCinema(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "14:00", result[0].ShowTime)
	assert.Equal(t, "Avengers", result[0].Movie.Title)
	assert.Equal(t, "Studio 1", result[0].Studio.Name)
	mockSeatRepo.AssertExpectations(t)
}

func TestSeatUseCase_GetShowtimesByCinema_Error(t *testing.T) {
	mockSeatRepo := new(MockSeatRepo)
	repo := &repository.Repository{Seat: mockSeatRepo}
	usecase := &SeatUseCase{Repo: repo}

	mockSeatRepo.On("GetShowtimesByCinema", mock.Anything, 999).Return([]entity.Showtime{}, errors.New("cinema not found"))

	result, err := usecase.GetShowtimesByCinema(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockSeatRepo.AssertExpectations(t)
}

func TestSeatUseCase_GetShowtimesByCinema_Empty(t *testing.T) {
	mockSeatRepo := new(MockSeatRepo)
	repo := &repository.Repository{Seat: mockSeatRepo}
	usecase := &SeatUseCase{Repo: repo}

	mockSeatRepo.On("GetShowtimesByCinema", mock.Anything, 1).Return([]entity.Showtime{}, nil)

	result, err := usecase.GetShowtimesByCinema(context.Background(), 1)

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockSeatRepo.AssertExpectations(t)
}

func TestSeatUseCase_GetShowtimesByCinema_WithNilMovieAndStudio(t *testing.T) {
	mockSeatRepo := new(MockSeatRepo)
	repo := &repository.Repository{Seat: mockSeatRepo}
	usecase := &SeatUseCase{Repo: repo}

	showtimes := []entity.Showtime{
		{ID: 1, CinemaID: 1, ShowDate: "2026-01-15", ShowTime: "14:00", Price: 50000, Movie: nil, Studio: nil},
	}

	mockSeatRepo.On("GetShowtimesByCinema", mock.Anything, 1).Return(showtimes, nil)

	result, err := usecase.GetShowtimesByCinema(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "14:00", result[0].ShowTime)
	// Movie and Studio should be zero values
	assert.Equal(t, 0, result[0].Movie.ID)
	assert.Equal(t, 0, result[0].Studio.ID)
	mockSeatRepo.AssertExpectations(t)
}

// =====================
// Mock Auth untuk Payment test dengan email
// =====================

type MockAuthRepoForPayment struct {
	mock.Mock
}

func (m *MockAuthRepoForPayment) CreateUser(ctx context.Context, user entity.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthRepoForPayment) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepoForPayment) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepoForPayment) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepoForPayment) UpdateUserVerified(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepoForPayment) CreateSession(ctx context.Context, session entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockAuthRepoForPayment) GetSessionByToken(ctx context.Context, token string) (entity.Session, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(entity.Session), args.Error(1)
}

func (m *MockAuthRepoForPayment) RevokeSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthRepoForPayment) CreateOTP(ctx context.Context, otp entity.OTP) error {
	args := m.Called(ctx, otp)
	return args.Error(0)
}

func (m *MockAuthRepoForPayment) GetValidOTP(ctx context.Context, userID int, otpCode string) (entity.OTP, error) {
	args := m.Called(ctx, userID, otpCode)
	return args.Get(0).(entity.OTP), args.Error(1)
}

func (m *MockAuthRepoForPayment) MarkOTPUsed(ctx context.Context, otpID int) error {
	args := m.Called(ctx, otpID)
	return args.Error(0)
}

func (m *MockAuthRepoForPayment) InvalidateUserOTPs(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestPaymentUseCase_ProcessPayment_Success(t *testing.T) {
	mockPaymentRepo := new(MockPaymentRepo)
	mockBookingRepo := new(MockBookingRepoForPayment)
	mockAuthRepo := new(MockAuthRepoForPayment)
	repo := &repository.Repository{
		Payment: mockPaymentRepo,
		Booking: mockBookingRepo,
		Auth:    mockAuthRepo,
	}
	usecase := &PaymentUseCase{Repo: repo, EmailService: utils.NewEmailService()}

	now := time.Now()
	booking := entity.Booking{
		ID:          1,
		UserID:      1,
		Status:      "pending",
		TotalAmount: 100000,
	}

	method := entity.PaymentMethod{
		ID:   1,
		Name: "Credit Card",
	}

	user := entity.User{
		ID:       1,
		Username: "testuser",
		Email:    "", // Empty email to skip goroutine
	}

	payment := entity.Payment{
		ID:        1,
		BookingID: 1,
		Status:    "completed",
		PaidAt:    &now,
	}

	mockBookingRepo.On("GetBookingByID", mock.Anything, 1).Return(booking, nil)
	mockPaymentRepo.On("GetPaymentMethodByID", mock.Anything, 1).Return(method, nil)
	mockPaymentRepo.On("CreatePayment", mock.Anything, mock.AnythingOfType("entity.Payment")).Return(1, nil)
	mockBookingRepo.On("UpdateBookingStatus", mock.Anything, 1, "paid").Return(nil)
	mockPaymentRepo.On("GetPaymentByBookingID", mock.Anything, 1).Return(payment, nil)
	mockAuthRepo.On("GetUserByID", mock.Anything, 1).Return(user, nil)

	req := dto.PayRequest{
		BookingID:     1,
		PaymentMethod: 1,
	}

	result, err := usecase.ProcessPayment(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, "Credit Card", result.PaymentMethod)
	mockBookingRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}
