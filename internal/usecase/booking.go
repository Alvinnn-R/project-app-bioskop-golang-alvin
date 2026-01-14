package usecase

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
)

type BookingUseCaseInterface interface {
	CreateBooking(ctx context.Context, userID int, req dto.BookingRequest) (dto.BookingResponse, error)
	GetUserBookings(ctx context.Context, userID int) ([]dto.BookingResponse, error)
}

type BookingUseCase struct {
	Repo *repository.Repository
}

func NewBookingUseCase(repo *repository.Repository) BookingUseCaseInterface {
	return &BookingUseCase{Repo: repo}
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

	return dto.BookingResponse{
		ID: bookingID,
		Showtime: dto.ShowtimeResponse{
			ID:       showtime.ID,
			ShowDate: showtime.ShowDate,
			ShowTime: showtime.ShowTime,
			Price:    showtime.Price,
		},
		Seats:       seatResponses,
		TotalAmount: createdBooking.TotalAmount,
		Status:      createdBooking.Status,
		CreatedAt:   createdBooking.CreatedAt,
	}, nil
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
