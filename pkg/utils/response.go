package utils

import (
	"encoding/json"
	"net/http"
	"session-23/internal/dto"
)

type Reponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func ResponseSuccess(w http.ResponseWriter, code int, message string, data any) {
	response := Reponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func ResponseBadRequest(w http.ResponseWriter, code int, message string, errors any) {
	response := Reponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func ResponsePagination(w http.ResponseWriter, code int, message string, data any, pagination dto.Pagination) {
	response := map[string]interface{}{
		"status":     true,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
