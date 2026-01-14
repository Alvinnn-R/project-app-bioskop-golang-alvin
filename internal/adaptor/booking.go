package adaptor

import (
	"encoding/json"
	"net/http"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type BookingAdaptor struct {
	UseCase  usecase.BookingUseCaseInterface
	Validate *validator.Validate
}

func NewBookingAdaptor(useCase usecase.BookingUseCaseInterface) *BookingAdaptor {
	return &BookingAdaptor{
		UseCase:  useCase,
		Validate: validator.New(),
	}
}

// Create handles seat booking creation
func (a *BookingAdaptor) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ResponseBadRequest(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := a.Validate.Struct(req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	booking, err := a.UseCase.CreateBooking(r.Context(), userID, req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "booking created successfully", booking)
}

// GetUserBookings handles get user booking history
func (a *BookingAdaptor) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ResponseBadRequest(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	bookings, err := a.UseCase.GetUserBookings(r.Context(), userID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to get bookings", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get booking history", bookings)
}
