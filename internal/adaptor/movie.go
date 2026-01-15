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
	Config  utils.Configuration
}

func NewMovieAdaptor(useCase usecase.MovieUseCaseInterface, config utils.Configuration) *MovieAdaptor {
	return &MovieAdaptor{
		UseCase: useCase,
		Config:  config,
	}
}

// GetAll handles get all movies with pagination
func (a *MovieAdaptor) GetAll(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get limit from config
	limit := a.Config.Limit
	if limit < 1 {
		limit = 10
	}

	movies, pagination, err := a.UseCase.GetAllMovies(r.Context(), page, limit)
	if err != nil {
		utils.ResponseInternalError(w, "failed to get movies")
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
		utils.ResponseNotFound(w, "movie not found")
		return
	}

	utils.ResponseOK(w, "success get movie", movie)
}
