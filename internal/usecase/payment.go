package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"project-app-bioskop/internal/data/entity"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/pkg/utils"
)

type PaymentUseCaseInterface interface {
	GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error)
	ProcessPayment(ctx context.Context, userID int, req dto.PayRequest) (dto.PaymentResponse, error)
}

type PaymentUseCase struct {
	Repo         *repository.Repository
	EmailService *utils.EmailService
}

func NewPaymentUseCase(repo *repository.Repository) PaymentUseCaseInterface {
	return &PaymentUseCase{
		Repo:         repo,
		EmailService: utils.NewEmailService(),
	}
}

// sendPaymentConfirmation sends payment confirmation email asynchronously (GOROUTINE)
func (u *PaymentUseCase) sendPaymentConfirmation(email, username, paymentMethod string, amount float64) {
	message := fmt.Sprintf(
		"Dear %s,\n\nYour payment has been processed successfully!\n\nPayment Method: %s\nAmount: Rp %.0f\n\nThank you for your purchase!",
		username, paymentMethod, amount,
	)
	u.EmailService.SendOTP(email, username, message)
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

	// Convert payment details to JSON format for jsonb column
	var paymentDetailsJSON []byte
	if req.PaymentDetails != "" {
		// Wrap string in JSON object if it's not already valid JSON
		detailsMap := map[string]string{"info": req.PaymentDetails}
		paymentDetailsJSON, _ = json.Marshal(detailsMap)
	} else {
		paymentDetailsJSON = []byte("{}")
	}

	// Create payment record
	payment := entity.Payment{
		BookingID:       req.BookingID,
		PaymentMethodID: req.PaymentMethod,
		Status:          "completed",
		PaymentDetails:  string(paymentDetailsJSON),
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

	// Send payment confirmation email asynchronously using GOROUTINE
	// This allows the API to respond immediately without waiting for email to be sent
	user, _ := u.Repo.Auth.GetUserByID(ctx, userID)
	if user.Email != "" {
		// GOROUTINE: Non-blocking email notification after payment
		go u.sendPaymentConfirmation(user.Email, user.Username, method.Name, booking.TotalAmount)
	}

	return dto.PaymentResponse{
		ID:            paymentID,
		PaymentMethod: method.Name,
		Status:        "completed",
		PaidAt:        createdPayment.PaidAt,
	}, nil
}
