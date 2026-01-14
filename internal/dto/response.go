package dto

import (
	"time"
)

// ResponseUser for user data response
type ResponseUser struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// LoginResponse for login success response
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiredAt time.Time    `json:"expired_at"`
	User      ResponseUser `json:"user"`
}

// CinemaResponse for cinema detail response
type CinemaResponse struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Location  string           `json:"location"`
	Studios   []StudioResponse `json:"studios,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

// StudioResponse for studio data response
type StudioResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	TotalSeats int    `json:"total_seats"`
}

// SeatResponse for seat availability response
type SeatResponse struct {
	ID       int    `json:"id"`
	SeatCode string `json:"seat_code"`
	IsBooked bool   `json:"is_booked"`
}

// ShowtimeResponse for showtime data response
type ShowtimeResponse struct {
	ID       int            `json:"id"`
	Movie    MovieResponse  `json:"movie"`
	Studio   StudioResponse `json:"studio"`
	ShowDate string         `json:"show_date"`
	ShowTime string         `json:"show_time"`
	Price    float64        `json:"price"`
}

// MovieResponse for movie data response
type MovieResponse struct {
	ID              int      `json:"id"`
	Title           string   `json:"title"`
	PosterURL       string   `json:"poster_url"`
	Genres          []string `json:"genres"`
	Rating          float64  `json:"rating"`
	DurationMinutes int      `json:"duration_minutes"`
}

// BookingResponse for booking data response
type BookingResponse struct {
	ID          int              `json:"id"`
	Showtime    ShowtimeResponse `json:"showtime"`
	Seats       []SeatResponse   `json:"seats"`
	TotalAmount float64          `json:"total_amount"`
	Status      string           `json:"status"`
	Payment     *PaymentResponse `json:"payment,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

// PaymentResponse for payment data response
type PaymentResponse struct {
	ID            int        `json:"id"`
	PaymentMethod string     `json:"payment_method"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}

// PaymentMethodResponse for payment method list response
type PaymentMethodResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PriceStats for dashboard statistics
type PriceStats struct {
	Min float64
	Max float64
	Avg float64
}
