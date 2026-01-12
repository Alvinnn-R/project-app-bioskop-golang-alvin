package wire

import (
	"session-23/internal/adaptor"
	"session-23/internal/data/repository"
	"session-23/internal/middleware"
	"session-23/internal/usecase"
	"session-23/pkg/utils"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Wiring(repo repository.Repository, config utils.Configuration, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	
	// Chi built-in middleware
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Recoverer)
	
	// Custom logging middleware to track duration
	mw := middleware.NewMiddlewareCustome(logger)
	router.Use(mw.Logging)
	
	// Mount routes
	router.Mount("/api/v1", ApiV1(repo, config))
	
	return router
}

func ApiV1(repo repository.Repository, config utils.Configuration) *chi.Mux {
	r := chi.NewRouter()
	
	// Wire car routes
	useCaseCar := usecase.NewServiceCar(&repo)
	adaptorCar := adaptor.NewAdaptorCar(useCaseCar, config)
	
	// Car dashboard routes
	r.Route("/cars", func(r chi.Router) {
		r.Get("/dashboard-serial", adaptorCar.Dashboard)
		r.Get("/dashboard-concurrent", adaptorCar.DashboardConcurrent)
	})
	
	return r
}
