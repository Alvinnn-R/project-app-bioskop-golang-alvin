package adaptor

import (
	"encoding/json"
	"net/http"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type AuthAdaptor struct {
	UseCase  usecase.AuthUseCaseInterface
	Validate *validator.Validate
}

func NewAuthAdaptor(useCase usecase.AuthUseCaseInterface) *AuthAdaptor {
	return &AuthAdaptor{
		UseCase:  useCase,
		Validate: validator.New(),
	}
}

// Register handles user registration
func (a *AuthAdaptor) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := a.Validate.Struct(req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	user, err := a.UseCase.Register(r.Context(), req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "registration successful", user)
}

// Login handles user login
func (a *AuthAdaptor) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := a.Validate.Struct(req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	response, err := a.UseCase.Login(r.Context(), req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "login successful", response)
}

// Logout handles user logout
func (a *AuthAdaptor) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "missing authorization token", nil)
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := a.UseCase.Logout(r.Context(), token); err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "logout failed", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "logout successful", nil)
}
