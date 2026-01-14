package adaptor

import (
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"
)

type Adaptor struct {
	AuthAdaptor    *AuthAdaptor
	CinemaAdaptor  *CinemaAdaptor
	SeatAdaptor    *SeatAdaptor
	BookingAdaptor *BookingAdaptor
	PaymentAdaptor *PaymentAdaptor
}

func NewAdaptor(repo *repository.Repository, config utils.Configuration) *Adaptor {
	// Initialize all usecases
	authUseCase := usecase.NewAuthUseCase(repo)
	cinemaUseCase := usecase.NewCinemaUseCase(repo)
	seatUseCase := usecase.NewSeatUseCase(repo)
	bookingUseCase := usecase.NewBookingUseCase(repo)
	paymentUseCase := usecase.NewPaymentUseCase(repo)

	return &Adaptor{
		AuthAdaptor:    NewAuthAdaptor(authUseCase),
		CinemaAdaptor:  NewCinemaAdaptor(cinemaUseCase),
		SeatAdaptor:    NewSeatAdaptor(seatUseCase),
		BookingAdaptor: NewBookingAdaptor(bookingUseCase),
		PaymentAdaptor: NewPaymentAdaptor(paymentUseCase),
	}
}
