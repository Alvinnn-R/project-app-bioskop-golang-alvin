package utils

import (
	"encoding/json"
	"net/http"
	"project-app-bioskop/internal/dto"
)

type Reponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// ResponseSuccess returns a successful response with custom status code
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

// ResponseCreated returns 201 Created response
func ResponseCreated(w http.ResponseWriter, message string, data any) {
	ResponseSuccess(w, http.StatusCreated, message, data)
}

// ResponseOK returns 200 OK response
func ResponseOK(w http.ResponseWriter, message string, data any) {
	ResponseSuccess(w, http.StatusOK, message, data)
}

// ResponseError returns an error response with custom status code
func ResponseError(w http.ResponseWriter, code int, message string, errors any) {
	response := Reponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// ResponseBadRequest returns 400 Bad Request response
func ResponseBadRequest(w http.ResponseWriter, code int, message string, errors any) {
	ResponseError(w, code, message, errors)
}

// ResponseUnauthorized returns 401 Unauthorized response
func ResponseUnauthorized(w http.ResponseWriter, message string) {
	ResponseError(w, http.StatusUnauthorized, message, nil)
}

// ResponseForbidden returns 403 Forbidden response
func ResponseForbidden(w http.ResponseWriter, message string) {
	ResponseError(w, http.StatusForbidden, message, nil)
}

// ResponseNotFound returns 404 Not Found response
func ResponseNotFound(w http.ResponseWriter, message string) {
	ResponseError(w, http.StatusNotFound, message, nil)
}

// ResponseInternalError returns 500 Internal Server Error response
func ResponseInternalError(w http.ResponseWriter, message string) {
	ResponseError(w, http.StatusInternalServerError, message, nil)
}

// ResponseValidationError returns 400 Bad Request with validation errors
func ResponseValidationError(w http.ResponseWriter, errors any) {
	ResponseError(w, http.StatusBadRequest, "validation error", errors)
}

// ResponsePagination returns paginated response
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
