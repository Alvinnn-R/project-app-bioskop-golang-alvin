package adaptor

import (
	"net/http"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MovieAdaptor struct {
	UseCase usecase.MovieUseCaseInterface
}

func NewMovieAdaptor(useCase usecase.MovieUseCaseInterface) *MovieAdaptor {
	return &MovieAdaptor{UseCase: useCase}
}

// GetAll handles get all movies with pagination
func (a *MovieAdaptor) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	movies, pagination, err := a.UseCase.GetAllMovies(r.Context(), page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to get movies", nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get movies", movies, pagination)
}

// GetByID handles get movie detail by ID
func (a *MovieAdaptor) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "movieId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid movie id", nil)
		return
	}

	movie, err := a.UseCase.GetMovieByID(r.Context(), id)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "movie not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get movie", movie)
}
