package adaptor

import (
	"net/http"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CinemaAdaptor struct {
	UseCase usecase.CinemaUseCaseInterface
	Config  utils.Configuration
}

func NewCinemaAdaptor(useCase usecase.CinemaUseCaseInterface, config utils.Configuration) *CinemaAdaptor {
	return &CinemaAdaptor{
		UseCase: useCase,
		Config:  config,
	}
}

// GetAll handles get all cinemas with pagination
func (a *CinemaAdaptor) GetAll(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get limit from config
	limit := a.Config.Limit
	if limit < 1 {
		limit = 10
	}

	cinemas, pagination, err := a.UseCase.GetAllCinemas(r.Context(), page, limit)
	if err != nil {
		utils.ResponseInternalError(w, "failed to get cinemas")
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get cinemas", cinemas, pagination)
}

// GetByID handles get cinema detail by ID
func (a *CinemaAdaptor) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cinemaId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid cinema id", nil)
		return
	}

	cinema, err := a.UseCase.GetCinemaByID(r.Context(), id)
	if err != nil {
		utils.ResponseNotFound(w, "cinema not found")
		return
	}

	utils.ResponseOK(w, "success get cinema", cinema)
}
