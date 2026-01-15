package entity

import "time"

// OTP represents email verification OTP
type OTP struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	OTPCode   string    `json:"otp_code"`
	ExpiredAt time.Time `json:"expired_at"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}
