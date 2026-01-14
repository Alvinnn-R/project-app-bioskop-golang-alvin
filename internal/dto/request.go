package dto

// RegisterRequest for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// BookingRequest for seat booking
type BookingRequest struct {
	ShowtimeID    int   `json:"showtime_id" validate:"required"`
	SeatIDs       []int `json:"seat_ids" validate:"required,min=1"`
	PaymentMethod int   `json:"payment_method" validate:"required"`
}

// PayRequest for processing payment
type PayRequest struct {
	BookingID      int    `json:"booking_id" validate:"required"`
	PaymentMethod  int    `json:"payment_method" validate:"required"`
	PaymentDetails string `json:"payment_details"`
}
