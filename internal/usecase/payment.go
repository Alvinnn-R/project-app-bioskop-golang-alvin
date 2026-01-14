package usecase

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
)

type PaymentUseCaseInterface interface {
	GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error)
	ProcessPayment(ctx context.Context, userID int, req dto.PayRequest) (dto.PaymentResponse, error)
}

type PaymentUseCase struct {
	Repo *repository.Repository
}

func NewPaymentUseCase(repo *repository.Repository) PaymentUseCaseInterface {
	return &PaymentUseCase{Repo: repo}
}

// GetPaymentMethods retrieves all available payment methods
func (u *PaymentUseCase) GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error) {
	methods, err := u.Repo.Payment.GetAllPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.PaymentMethodResponse
	for _, m := range methods {
		response = append(response, dto.PaymentMethodResponse{
			ID:   m.ID,
			Name: m.Name,
		})
	}

	return response, nil
}

// ProcessPayment processes payment for a booking
func (u *PaymentUseCase) ProcessPayment(ctx context.Context, userID int, req dto.PayRequest) (dto.PaymentResponse, error) {
	// Verify booking exists and belongs to user
	booking, err := u.Repo.Booking.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return dto.PaymentResponse{}, errors.New("booking not found")
	}

	if booking.UserID != userID {
		return dto.PaymentResponse{}, errors.New("unauthorized to pay for this booking")
	}

	if booking.Status == "paid" {
		return dto.PaymentResponse{}, errors.New("booking already paid")
	}

	if booking.Status == "cancelled" {
		return dto.PaymentResponse{}, errors.New("booking is cancelled")
	}

	// Verify payment method exists
	method, err := u.Repo.Payment.GetPaymentMethodByID(ctx, req.PaymentMethod)
	if err != nil {
		return dto.PaymentResponse{}, errors.New("invalid payment method")
	}

	// Create payment record
	payment := entity.Payment{
		BookingID:       req.BookingID,
		PaymentMethodID: req.PaymentMethod,
		Status:          "completed",
		PaymentDetails:  req.PaymentDetails,
	}

	paymentID, err := u.Repo.Payment.CreatePayment(ctx, payment)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	// Update booking status to paid
	if err := u.Repo.Booking.UpdateBookingStatus(ctx, req.BookingID, "paid"); err != nil {
		return dto.PaymentResponse{}, err
	}

	createdPayment, _ := u.Repo.Payment.GetPaymentByBookingID(ctx, req.BookingID)

	return dto.PaymentResponse{
		ID:            paymentID,
		PaymentMethod: method.Name,
		Status:        "completed",
		PaidAt:        createdPayment.PaidAt,
	}, nil
}
