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
// Mock Repository untuk Booking
// =====================

type MockBookingRepo struct {
	mock.Mock
}

func (m *MockBookingRepo) CreateBooking(ctx context.Context, booking entity.Booking, seatIDs []int, price float64) (int, error) {
	args := m.Called(ctx, booking, seatIDs, price)
	return args.Int(0), args.Error(1)
}

func (m *MockBookingRepo) GetBookingByID(ctx context.Context, id int) (entity.Booking, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Booking), args.Error(1)
}

func (m *MockBookingRepo) GetBookingsByUserID(ctx context.Context, userID int) ([]entity.Booking, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Booking), args.Error(1)
}

func (m *MockBookingRepo) GetBookingSeats(ctx context.Context, bookingID int) ([]entity.BookingSeat, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).([]entity.BookingSeat), args.Error(1)
}

func (m *MockBookingRepo) UpdateBookingStatus(ctx context.Context, bookingID int, status string) error {
	args := m.Called(ctx, bookingID, status)
	return args.Error(0)
}

// =====================
// Mock Repository untuk Seat (Booking test)
// =====================

type MockSeatRepoForBooking struct {
	mock.Mock
}

func (m *MockSeatRepoForBooking) GetSeatsByShowtime(ctx context.Context, showtimeID int) ([]entity.SeatAvailability, error) {
	args := m.Called(ctx, showtimeID)
	return args.Get(0).([]entity.SeatAvailability), args.Error(1)
}

func (m *MockSeatRepoForBooking) GetShowtimeByParams(ctx context.Context, cinemaID int, date, time string) (entity.Showtime, error) {
	args := m.Called(ctx, cinemaID, date, time)
	return args.Get(0).(entity.Showtime), args.Error(1)
}

func (m *MockSeatRepoForBooking) GetShowtimeByID(ctx context.Context, id int) (entity.Showtime, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Showtime), args.Error(1)
}

func (m *MockSeatRepoForBooking) GetSeatsByIDs(ctx context.Context, seatIDs []int) ([]entity.Seat, error) {
	args := m.Called(ctx, seatIDs)
	return args.Get(0).([]entity.Seat), args.Error(1)
}

func (m *MockSeatRepoForBooking) CheckSeatsAvailable(ctx context.Context, showtimeID int, seatIDs []int) (bool, error) {
	args := m.Called(ctx, showtimeID, seatIDs)
	return args.Bool(0), args.Error(1)
}

func (m *MockSeatRepoForBooking) GetShowtimesByCinema(ctx context.Context, cinemaID int) ([]entity.Showtime, error) {
	args := m.Called(ctx, cinemaID)
	return args.Get(0).([]entity.Showtime), args.Error(1)
}

// =====================
// Mock Payment Repo (Booking test)
// =====================

type MockPaymentRepoForBooking struct {
	mock.Mock
}

func (m *MockPaymentRepoForBooking) GetAllPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.PaymentMethod), args.Error(1)
}

func (m *MockPaymentRepoForBooking) GetPaymentMethodByID(ctx context.Context, id int) (entity.PaymentMethod, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.PaymentMethod), args.Error(1)
}

func (m *MockPaymentRepoForBooking) CreatePayment(ctx context.Context, payment entity.Payment) (int, error) {
	args := m.Called(ctx, payment)
	return args.Int(0), args.Error(1)
}

func (m *MockPaymentRepoForBooking) GetPaymentByBookingID(ctx context.Context, bookingID int) (entity.Payment, error) {
	args := m.Called(ctx, bookingID)
	return args.Get(0).(entity.Payment), args.Error(1)
}

func (m *MockPaymentRepoForBooking) UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error {
	args := m.Called(ctx, paymentID, status)
	return args.Error(0)
}

// =====================
// Mock Auth Repo (Booking test)
// =====================

type MockAuthRepoForBooking struct {
	mock.Mock
}

func (m *MockAuthRepoForBooking) CreateUser(ctx context.Context, user entity.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthRepoForBooking) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepoForBooking) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepoForBooking) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthRepoForBooking) UpdateUserVerified(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepoForBooking) CreateSession(ctx context.Context, session entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockAuthRepoForBooking) GetSessionByToken(ctx context.Context, token string) (entity.Session, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(entity.Session), args.Error(1)
}

func (m *MockAuthRepoForBooking) RevokeSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthRepoForBooking) CreateOTP(ctx context.Context, otp entity.OTP) error {
	args := m.Called(ctx, otp)
	return args.Error(0)
}

func (m *MockAuthRepoForBooking) GetValidOTP(ctx context.Context, userID int, otpCode string) (entity.OTP, error) {
	args := m.Called(ctx, userID, otpCode)
	return args.Get(0).(entity.OTP), args.Error(1)
}

func (m *MockAuthRepoForBooking) MarkOTPUsed(ctx context.Context, otpID int) error {
	args := m.Called(ctx, otpID)
	return args.Error(0)
}

func (m *MockAuthRepoForBooking) InvalidateUserOTPs(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// =====================
// Booking UseCase Tests
// =====================

func TestBookingUseCase_CreateBooking_ShowtimeNotFound(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockSeatRepo := new(MockSeatRepoForBooking)
	mockPaymentRepo := new(MockPaymentRepoForBooking)

	repo := &repository.Repository{
		Booking: mockBookingRepo,
		Seat:    mockSeatRepo,
		Payment: mockPaymentRepo,
	}
	usecase := &BookingUseCase{Repo: repo}

	mockSeatRepo.On("GetShowtimeByID", mock.Anything, 999).Return(entity.Showtime{}, errors.New("not found"))

	req := dto.BookingRequest{
		ShowtimeID:    999,
		SeatIDs:       []int{1, 2},
		PaymentMethod: 1,
	}

	_, err := usecase.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "showtime not found")
	mockSeatRepo.AssertExpectations(t)
}

func TestBookingUseCase_CreateBooking_SeatsNotAvailable(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockSeatRepo := new(MockSeatRepoForBooking)
	mockPaymentRepo := new(MockPaymentRepoForBooking)

	repo := &repository.Repository{
		Booking: mockBookingRepo,
		Seat:    mockSeatRepo,
		Payment: mockPaymentRepo,
	}
	usecase := &BookingUseCase{Repo: repo}

	movie := &entity.Movie{ID: 1, Title: "Avengers", Genres: []string{"Action"}}
	studio := &entity.Studio{ID: 1, Name: "Studio 1", TotalSeats: 100}
	showtime := entity.Showtime{
		ID:       1,
		CinemaID: 1,
		ShowDate: "2026-01-15",
		ShowTime: "19:00",
		Price:    50000,
		Movie:    movie,
		Studio:   studio,
	}

	mockSeatRepo.On("GetShowtimeByID", mock.Anything, 1).Return(showtime, nil)
	mockSeatRepo.On("CheckSeatsAvailable", mock.Anything, 1, []int{1, 2}).Return(false, nil)

	req := dto.BookingRequest{
		ShowtimeID:    1,
		SeatIDs:       []int{1, 2},
		PaymentMethod: 1,
	}

	_, err := usecase.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not available")
	mockSeatRepo.AssertExpectations(t)
}

func TestBookingUseCase_CreateBooking_InvalidPaymentMethod(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockSeatRepo := new(MockSeatRepoForBooking)
	mockPaymentRepo := new(MockPaymentRepoForBooking)

	repo := &repository.Repository{
		Booking: mockBookingRepo,
		Seat:    mockSeatRepo,
		Payment: mockPaymentRepo,
	}
	usecase := &BookingUseCase{Repo: repo}

	movie := &entity.Movie{ID: 1, Title: "Avengers", Genres: []string{"Action"}}
	studio := &entity.Studio{ID: 1, Name: "Studio 1", TotalSeats: 100}
	showtime := entity.Showtime{
		ID:       1,
		CinemaID: 1,
		ShowDate: "2026-01-15",
		ShowTime: "19:00",
		Price:    50000,
		Movie:    movie,
		Studio:   studio,
	}

	mockSeatRepo.On("GetShowtimeByID", mock.Anything, 1).Return(showtime, nil)
	mockSeatRepo.On("CheckSeatsAvailable", mock.Anything, 1, []int{1, 2}).Return(true, nil)
	mockPaymentRepo.On("GetPaymentMethodByID", mock.Anything, 999).Return(entity.PaymentMethod{}, errors.New("not found"))

	req := dto.BookingRequest{
		ShowtimeID:    1,
		SeatIDs:       []int{1, 2},
		PaymentMethod: 999,
	}

	_, err := usecase.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payment method")
	mockSeatRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
}

func TestBookingUseCase_CreateBooking_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockSeatRepo := new(MockSeatRepoForBooking)
	mockPaymentRepo := new(MockPaymentRepoForBooking)
	mockAuthRepo := new(MockAuthRepoForBooking)

	repo := &repository.Repository{
		Booking: mockBookingRepo,
		Seat:    mockSeatRepo,
		Payment: mockPaymentRepo,
		Auth:    mockAuthRepo,
	}
	usecase := &BookingUseCase{Repo: repo, EmailService: utils.NewEmailService()}

	now := time.Now()
	movie := &entity.Movie{ID: 1, Title: "Avengers", PosterURL: "http://poster.url", Genres: []string{"Action"}, Rating: 8.5, DurationMinutes: 180}
	studio := &entity.Studio{ID: 1, Name: "Studio 1", TotalSeats: 100}
	showtime := entity.Showtime{
		ID:       1,
		CinemaID: 1,
		ShowDate: "2026-01-15",
		ShowTime: "19:00",
		Price:    50000,
		Movie:    movie,
		Studio:   studio,
	}

	seats := []entity.Seat{
		{ID: 1, SeatCode: "A1", StudioID: 1},
		{ID: 2, SeatCode: "A2", StudioID: 1},
	}

	createdBooking := entity.Booking{
		ID:          1,
		UserID:      1,
		ShowtimeID:  1,
		Status:      "pending",
		TotalAmount: 100000,
		CreatedAt:   now,
	}

	user := entity.User{
		ID:       1,
		Username: "testuser",
		Email:    "", // Empty email to skip goroutine
	}

	mockSeatRepo.On("GetShowtimeByID", mock.Anything, 1).Return(showtime, nil)
	mockSeatRepo.On("CheckSeatsAvailable", mock.Anything, 1, []int{1, 2}).Return(true, nil)
	mockPaymentRepo.On("GetPaymentMethodByID", mock.Anything, 1).Return(entity.PaymentMethod{ID: 1, Name: "Credit Card"}, nil)
	mockBookingRepo.On("CreateBooking", mock.Anything, mock.AnythingOfType("entity.Booking"), []int{1, 2}, 50000.0).Return(1, nil)
	mockBookingRepo.On("GetBookingByID", mock.Anything, 1).Return(createdBooking, nil)
	mockSeatRepo.On("GetSeatsByIDs", mock.Anything, []int{1, 2}).Return(seats, nil)
	mockAuthRepo.On("GetUserByID", mock.Anything, 1).Return(user, nil)

	req := dto.BookingRequest{
		ShowtimeID:    1,
		SeatIDs:       []int{1, 2},
		PaymentMethod: 1,
	}

	result, err := usecase.CreateBooking(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, 100000.0, result.TotalAmount)
	assert.Len(t, result.Seats, 2)
	assert.Equal(t, "Avengers", result.Showtime.Movie.Title)
	mockSeatRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_GetUserBookings_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockSeatRepo := new(MockSeatRepoForBooking)
	mockPaymentRepo := new(MockPaymentRepoForBooking)

	repo := &repository.Repository{
		Booking: mockBookingRepo,
		Seat:    mockSeatRepo,
		Payment: mockPaymentRepo,
	}
	usecase := &BookingUseCase{Repo: repo}

	now := time.Now()
	movie := &entity.Movie{ID: 1, Title: "Avengers", Genres: []string{"Action"}, Rating: 8.5, DurationMinutes: 180}
	studio := &entity.Studio{ID: 1, Name: "Studio 1", TotalSeats: 100}

	bookings := []entity.Booking{
		{ID: 1, UserID: 1, ShowtimeID: 1, Status: "paid", TotalAmount: 100000, CreatedAt: now},
	}

	showtime1 := entity.Showtime{ID: 1, ShowDate: "2026-01-15", ShowTime: "19:00", Price: 50000, Movie: movie, Studio: studio}

	bookingSeats := []entity.BookingSeat{
		{ID: 1, BookingID: 1, SeatID: 1, PriceSnapshot: 50000},
		{ID: 2, BookingID: 1, SeatID: 2, PriceSnapshot: 50000},
	}

	seats := []entity.Seat{
		{ID: 1, SeatCode: "A1", StudioID: 1},
		{ID: 2, SeatCode: "A2", StudioID: 1},
	}

	payment := entity.Payment{
		ID:              1,
		BookingID:       1,
		PaymentMethodID: 1,
		Status:          "completed",
		PaidAt:          &now,
	}

	paymentMethod := entity.PaymentMethod{ID: 1, Name: "Credit Card"}

	mockBookingRepo.On("GetBookingsByUserID", mock.Anything, 1).Return(bookings, nil)
	mockSeatRepo.On("GetShowtimeByID", mock.Anything, 1).Return(showtime1, nil)
	mockBookingRepo.On("GetBookingSeats", mock.Anything, 1).Return(bookingSeats, nil)
	mockSeatRepo.On("GetSeatsByIDs", mock.Anything, []int{1, 2}).Return(seats, nil)
	mockPaymentRepo.On("GetPaymentByBookingID", mock.Anything, 1).Return(payment, nil)
	mockPaymentRepo.On("GetPaymentMethodByID", mock.Anything, 1).Return(paymentMethod, nil)

	result, err := usecase.GetUserBookings(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "paid", result[0].Status)
	mockBookingRepo.AssertExpectations(t)
	mockSeatRepo.AssertExpectations(t)
}

func TestBookingUseCase_GetUserBookings_Error(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockSeatRepo := new(MockSeatRepoForBooking)

	repo := &repository.Repository{
		Booking: mockBookingRepo,
		Seat:    mockSeatRepo,
	}
	usecase := &BookingUseCase{Repo: repo}

	mockBookingRepo.On("GetBookingsByUserID", mock.Anything, 1).Return([]entity.Booking{}, errors.New("database error"))

	result, err := usecase.GetUserBookings(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_CreateBooking_CreateBookingError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockSeatRepo := new(MockSeatRepoForBooking)
	mockPaymentRepo := new(MockPaymentRepoForBooking)

	repo := &repository.Repository{
		Booking: mockBookingRepo,
		Seat:    mockSeatRepo,
		Payment: mockPaymentRepo,
	}
	usecase := &BookingUseCase{Repo: repo}

	movie := &entity.Movie{ID: 1, Title: "Avengers", Genres: []string{"Action"}}
	studio := &entity.Studio{ID: 1, Name: "Studio 1", TotalSeats: 100}
	showtime := entity.Showtime{
		ID:       1,
		CinemaID: 1,
		ShowDate: "2026-01-15",
		ShowTime: "19:00",
		Price:    50000,
		Movie:    movie,
		Studio:   studio,
	}

	mockSeatRepo.On("GetShowtimeByID", mock.Anything, 1).Return(showtime, nil)
	mockSeatRepo.On("CheckSeatsAvailable", mock.Anything, 1, []int{1, 2}).Return(true, nil)
	mockPaymentRepo.On("GetPaymentMethodByID", mock.Anything, 1).Return(entity.PaymentMethod{ID: 1, Name: "Credit Card"}, nil)
	mockBookingRepo.On("CreateBooking", mock.Anything, mock.AnythingOfType("entity.Booking"), []int{1, 2}, 50000.0).Return(0, errors.New("database error"))

	req := dto.BookingRequest{
		ShowtimeID:    1,
		SeatIDs:       []int{1, 2},
		PaymentMethod: 1,
	}

	_, err := usecase.CreateBooking(context.Background(), 1, req)

	assert.Error(t, err)
	mockSeatRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
	mockBookingRepo.AssertExpectations(t)
}
