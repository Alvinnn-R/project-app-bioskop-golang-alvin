package adaptor

import (
	"net/http"
	"session-23/internal/usecase"
	"session-23/pkg/utils"
)

type AdaptorCar struct {
	Usecase usecase.ServiceCar
	Config  utils.Configuration
}

func NewAdaptorCar(usecase *usecase.ServiceCar, config utils.Configuration) *AdaptorCar {
	return &AdaptorCar{Usecase: *usecase, Config: config}
}

func (usecaseAdaptor *AdaptorCar) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get limit from config
	limit := usecaseAdaptor.Config.Limit

	// Get data from usecase with context
	response, err := usecaseAdaptor.Usecase.DashboardSerial(r.Context(), limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch cars: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data", response)
}

func (usecaseAdaptor *AdaptorCar) DashboardConcurrent(w http.ResponseWriter, r *http.Request) {
	// Get limit from config
	limit := usecaseAdaptor.Config.Limit

	// Get data from usecase with context using concurrent queries
	response, err := usecaseAdaptor.Usecase.DashboardConcurrent(r.Context(), limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch cars: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data with concurrent", response)
}
