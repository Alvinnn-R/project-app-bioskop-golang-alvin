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
}

func NewCinemaAdaptor(useCase usecase.CinemaUseCaseInterface) *CinemaAdaptor {
	return &CinemaAdaptor{UseCase: useCase}
}

// GetAll handles get all cinemas with pagination
func (a *CinemaAdaptor) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	cinemas, pagination, err := a.UseCase.GetAllCinemas(r.Context(), page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to get cinemas", nil)
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
		utils.ResponseBadRequest(w, http.StatusNotFound, "cinema not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get cinema", cinema)
}
