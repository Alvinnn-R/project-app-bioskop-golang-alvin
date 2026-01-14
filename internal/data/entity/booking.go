package entity

import "time"

// Booking represents a ticket booking
type Booking struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ShowtimeID  int       `json:"showtime_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

// BookingSeat represents a booked seat in a booking
type BookingSeat struct {
	ID            int     `json:"id"`
	BookingID     int     `json:"booking_id"`
	SeatID        int     `json:"seat_id"`
	PriceSnapshot float64 `json:"price_snapshot"`
}

// PaymentMethod represents available payment methods
type PaymentMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Payment represents a payment record
type Payment struct {
	ID              int        `json:"id"`
	BookingID       int        `json:"booking_id"`
	PaymentMethodID int        `json:"payment_method_id"`
	Status          string     `json:"status"`
	PaymentDetails  string     `json:"payment_details"`
	PaidAt          *time.Time `json:"paid_at"`
}
