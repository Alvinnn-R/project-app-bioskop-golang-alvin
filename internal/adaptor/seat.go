package adaptor

import (
	"net/http"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SeatAdaptor struct {
	UseCase usecase.SeatUseCaseInterface
}

func NewSeatAdaptor(useCase usecase.SeatUseCaseInterface) *SeatAdaptor {
	return &SeatAdaptor{UseCase: useCase}
}

// GetAvailability handles get seat availability for a cinema showtime
func (a *SeatAdaptor) GetAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cinemaId")
	cinemaID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid cinema id", nil)
		return
	}

	date := r.URL.Query().Get("date")
	showTime := r.URL.Query().Get("time")

	if date == "" || showTime == "" {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "date and time are required", nil)
		return
	}

	seats, err := a.UseCase.GetSeatAvailability(r.Context(), cinemaID, date, showTime)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "no showtime found for given parameters", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get seat availability", seats)
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
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to get showtimes", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get showtimes", showtimes)
}
