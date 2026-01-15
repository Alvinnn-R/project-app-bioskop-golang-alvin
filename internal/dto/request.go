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

// VerifyOTPRequest for email OTP verification
type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

// ResendOTPRequest for resending OTP
type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// BookingRequest for seat booking
type BookingRequest struct {
	ShowtimeID    int    `json:"showtime_id" validate:"required"`
	SeatIDs       []int  `json:"seat_ids" validate:"required,min=1"`
	PaymentMethod int    `json:"payment_method" validate:"required"`
	Date          string `json:"date" validate:"required"`
	Time          string `json:"time" validate:"required"`
}

// PayRequest for processing payment
type PayRequest struct {
	BookingID      int    `json:"booking_id" validate:"required"`
	PaymentMethod  int    `json:"payment_method" validate:"required"`
	PaymentDetails string `json:"payment_details"`
}

// SeatQueryRequest for seat availability query
type SeatQueryRequest struct {
	Date string `validate:"required"`
	Time string `validate:"required"`
}
