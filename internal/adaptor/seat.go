package adaptor

import (
	"net/http"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type SeatAdaptor struct {
	UseCase  usecase.SeatUseCaseInterface
	Validate *validator.Validate
}

func NewSeatAdaptor(useCase usecase.SeatUseCaseInterface) *SeatAdaptor {
	return &SeatAdaptor{
		UseCase:  useCase,
		Validate: validator.New(),
	}
}

// GetAvailability handles get seat availability for a cinema showtime
func (a *SeatAdaptor) GetAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cinemaId")
	cinemaID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid cinema id", nil)
		return
	}

	// Get query params
	req := dto.SeatQueryRequest{
		Date: r.URL.Query().Get("date"),
		Time: r.URL.Query().Get("time"),
	}

	// Validate query params
	if err := a.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, err.Error())
		return
	}

	seats, err := a.UseCase.GetSeatAvailability(r.Context(), cinemaID, req.Date, req.Time)
	if err != nil {
		utils.ResponseNotFound(w, "no showtime found for given parameters")
		return
	}

	utils.ResponseOK(w, "success get seat availability", seats)
}

// GetShowtimes handles get all showtimes for a cinema
func (a *SeatAdaptor) GetShowtimes(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cinemaId")
	cinemaID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid cinema id", nil)
		return
	}

	showtimes, err := a.UseCase.GetShowtimesByCinema(r.Context(), cinemaID)
	if err != nil {
		utils.ResponseInternalError(w, "failed to get showtimes")
		return
	}

	utils.ResponseOK(w, "success get showtimes", showtimes)
}
