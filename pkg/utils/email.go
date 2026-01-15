package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// EmailRequest represents the request body for email API
type EmailRequest struct {
	To      string `json:"to"`
	Name    string `json:"Name"`
	Subject string `json:"subject"`
	Message string `json:"body"`
}

// EmailService handles email sending operations
type EmailService struct {
	APIUrl string
	APIKey string
}

// NewEmailService creates a new EmailService instance
func NewEmailService() *EmailService {
	return &EmailService{
		APIUrl: "https://lumoshive-academy-email-api.vercel.app/send-email",
		APIKey: "0540540540540540540540540540540540540540540540540540540540540540",
	}
}

// SendOTP sends OTP code to user's email
func (e *EmailService) SendOTP(toEmail, name, otpCode string) error {
	payload := EmailRequest{
		To:      toEmail,
		Name:    name,
		Subject: "Email Verification - Cinema Booking System",
		Message: "Your OTP verification code is: " + otpCode + ". This code will expire in 5 minutes.",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.APIUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Setup headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", e.APIKey)

	// Send request
	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return errors.New("failed to send email")
	}

	return nil
}
