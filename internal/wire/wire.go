package wire

import (
	"project-app-bioskop/internal/adaptor"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/middleware"
	"project-app-bioskop/pkg/utils"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Wiring(repo *repository.Repository, config utils.Configuration, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	// Chi built-in middleware
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Recoverer)

	// Custom logging middleware
	mw := middleware.NewMiddlewareCustome(logger)
	router.Use(mw.Logging)

	// Initialize all adaptors
	adaptors := adaptor.NewAdaptor(repo, config)

	// Create auth usecase for middleware
	authUseCase := usecase.NewAuthUseCase(repo)
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// Mount API routes
	router.Route("/api", func(r chi.Router) {
		// Public routes - Authentication
		r.Post("/register", adaptors.AuthAdaptor.Register)
		r.Post("/verify-otp", adaptors.AuthAdaptor.VerifyOTP)
		r.Post("/resend-otp", adaptors.AuthAdaptor.ResendOTP)
		r.Post("/login", adaptors.AuthAdaptor.Login)

		// Public routes - Movies
		r.Get("/movies", adaptors.MovieAdaptor.GetAll)
		r.Get("/movies/{movieId}", adaptors.MovieAdaptor.GetByID)

		// Public routes - Cinema
		r.Get("/cinemas", adaptors.CinemaAdaptor.GetAll)
		r.Get("/cinemas/{cinemaId}", adaptors.CinemaAdaptor.GetByID)
		r.Get("/cinemas/{cinemaId}/showtimes", adaptors.SeatAdaptor.GetShowtimes)
		r.Get("/cinemas/{cinemaId}/seats", adaptors.SeatAdaptor.GetAvailability)

		// Public routes - Payment Methods
		r.Get("/payment-methods", adaptors.PaymentAdaptor.GetMethods)

		// Protected routes - require authentication
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			// Auth
			r.Post("/logout", adaptors.AuthAdaptor.Logout)

			// Booking
			r.Post("/booking", adaptors.BookingAdaptor.Create)
			r.Post("/pay", adaptors.PaymentAdaptor.ProcessPayment)

			// User routes
			r.Get("/user/bookings", adaptors.BookingAdaptor.GetUserBookings)
		})
	})

	return router
}
