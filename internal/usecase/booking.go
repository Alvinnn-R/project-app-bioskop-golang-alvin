package usecase

import (
	"context"
	"errors"
	"fmt"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/pkg/utils"
	"strings"
)

type BookingUseCaseInterface interface {
	CreateBooking(ctx context.Context, userID int, req dto.BookingRequest) (dto.BookingResponse, error)
	GetUserBookings(ctx context.Context, userID int) ([]dto.BookingResponse, error)
}

type BookingUseCase struct {
	Repo         *repository.Repository
	EmailService *utils.EmailService
}

func NewBookingUseCase(repo *repository.Repository) BookingUseCaseInterface {
	return &BookingUseCase{
		Repo:         repo,
		EmailService: utils.NewEmailService(),
	}
}

// sendBookingConfirmation sends booking confirmation email asynchronously (GOROUTINE)
func (u *BookingUseCase) sendBookingConfirmation(email, username, movieTitle, showDate, showTime string, seats []string, totalAmount float64) {
	seatList := strings.Join(seats, ", ")
	message := fmt.Sprintf(
		"Dear %s,\n\nYour booking has been confirmed!\n\nMovie: %s\nDate: %s\nTime: %s\nSeats: %s\nTotal Amount: Rp %.0f\n\nThank you for choosing our cinema!",
		username, movieTitle, showDate, showTime, seatList, totalAmount,
	)
	u.EmailService.SendOTP(email, username, message)
}

// CreateBooking creates a new seat booking
func (u *BookingUseCase) CreateBooking(ctx context.Context, userID int, req dto.BookingRequest) (dto.BookingResponse, error) {
	// Verify showtime exists
	showtime, err := u.Repo.Seat.GetShowtimeByID(ctx, req.ShowtimeID)
	if err != nil {
		return dto.BookingResponse{}, errors.New("showtime not found")
	}

	// Check seat availability
	available, err := u.Repo.Seat.CheckSeatsAvailable(ctx, req.ShowtimeID, req.SeatIDs)
	if err != nil {
		return dto.BookingResponse{}, err
	}
	if !available {
		return dto.BookingResponse{}, errors.New("one or more seats are not available")
	}

	// Verify payment method exists
	_, err = u.Repo.Payment.GetPaymentMethodByID(ctx, req.PaymentMethod)
	if err != nil {
		return dto.BookingResponse{}, errors.New("invalid payment method")
	}

	// Create booking
	booking := entity.Booking{
		UserID:     userID,
		ShowtimeID: req.ShowtimeID,
	}

	bookingID, err := u.Repo.Booking.CreateBooking(ctx, booking, req.SeatIDs, showtime.Price)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	// Get created booking for response
	createdBooking, err := u.Repo.Booking.GetBookingByID(ctx, bookingID)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	// Build response
	seats, _ := u.Repo.Seat.GetSeatsByIDs(ctx, req.SeatIDs)
	var seatResponses []dto.SeatResponse
	for _, s := range seats {
		seatResponses = append(seatResponses, dto.SeatResponse{
			ID:       s.ID,
			SeatCode: s.SeatCode,
			IsBooked: true,
		})
	}

	response := dto.BookingResponse{
		ID: bookingID,
		Showtime: dto.ShowtimeResponse{
			ID:       showtime.ID,
			ShowDate: showtime.ShowDate,
			ShowTime: showtime.ShowTime,
			Price:    showtime.Price,
			Movie: dto.MovieResponse{
				ID:              showtime.Movie.ID,
				Title:           showtime.Movie.Title,
				PosterURL:       showtime.Movie.PosterURL,
				Genres:          showtime.Movie.Genres,
				Rating:          showtime.Movie.Rating,
				DurationMinutes: showtime.Movie.DurationMinutes,
			},
			Studio: dto.StudioResponse{
				ID:         showtime.Studio.ID,
				Name:       showtime.Studio.Name,
				TotalSeats: showtime.Studio.TotalSeats,
			},
		},
		Seats:       seatResponses,
		TotalAmount: createdBooking.TotalAmount,
		Status:      createdBooking.Status,
		CreatedAt:   createdBooking.CreatedAt,
	}

	// Send booking confirmation email asynchronously using GOROUTINE
	// This allows the API to respond immediately without waiting for email to be sent
	user, _ := u.Repo.Auth.GetUserByID(ctx, userID)
	if user.Email != "" {
		var seatCodes []string
		for _, s := range seats {
			seatCodes = append(seatCodes, s.SeatCode)
		}
		// GOROUTINE: Non-blocking email notification
		go u.sendBookingConfirmation(
			user.Email,
			user.Username,
			showtime.Movie.Title,
			showtime.ShowDate,
			showtime.ShowTime,
			seatCodes,
			createdBooking.TotalAmount,
		)
	}

	return response, nil
}

// GetUserBookings retrieves all bookings for a user
func (u *BookingUseCase) GetUserBookings(ctx context.Context, userID int) ([]dto.BookingResponse, error) {
	bookings, err := u.Repo.Booking.GetBookingsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.BookingResponse
	for _, b := range bookings {
		showtime, _ := u.Repo.Seat.GetShowtimeByID(ctx, b.ShowtimeID)
		bookingSeats, _ := u.Repo.Booking.GetBookingSeats(ctx, b.ID)

		var seatIDs []int
		for _, bs := range bookingSeats {
			seatIDs = append(seatIDs, bs.SeatID)
		}

		seats, _ := u.Repo.Seat.GetSeatsByIDs(ctx, seatIDs)
		var seatResponses []dto.SeatResponse
		for _, s := range seats {
			seatResponses = append(seatResponses, dto.SeatResponse{
				ID:       s.ID,
				SeatCode: s.SeatCode,
			})
		}

		// Get payment if exists
		var paymentResp *dto.PaymentResponse
		payment, err := u.Repo.Payment.GetPaymentByBookingID(ctx, b.ID)
		if err == nil {
			method, _ := u.Repo.Payment.GetPaymentMethodByID(ctx, payment.PaymentMethodID)
			paymentResp = &dto.PaymentResponse{
				ID:            payment.ID,
				PaymentMethod: method.Name,
				Status:        payment.Status,
				PaidAt:        payment.PaidAt,
			}
		}

		responses = append(responses, dto.BookingResponse{
			ID: b.ID,
			Showtime: dto.ShowtimeResponse{
				ID:       showtime.ID,
				ShowDate: showtime.ShowDate,
				ShowTime: showtime.ShowTime,
				Price:    showtime.Price,
				Movie: dto.MovieResponse{
					ID:              showtime.Movie.ID,
					Title:           showtime.Movie.Title,
					PosterURL:       showtime.Movie.PosterURL,
					Genres:          showtime.Movie.Genres,
					Rating:          showtime.Movie.Rating,
					DurationMinutes: showtime.Movie.DurationMinutes,
				},
				Studio: dto.StudioResponse{
					ID:         showtime.Studio.ID,
					Name:       showtime.Studio.Name,
					TotalSeats: showtime.Studio.TotalSeats,
				},
			},
			Seats:       seatResponses,
			TotalAmount: b.TotalAmount,
			Status:      b.Status,
			Payment:     paymentResp,
			CreatedAt:   b.CreatedAt,
		})
	}

	return responses, nil
}
