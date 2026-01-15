package adaptor

import (
	"encoding/json"
	"net/http"
	"project-app-bioskop/internal/dto"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type PaymentAdaptor struct {
	UseCase  usecase.PaymentUseCaseInterface
	Validate *validator.Validate
}

func NewPaymentAdaptor(useCase usecase.PaymentUseCaseInterface) *PaymentAdaptor {
	return &PaymentAdaptor{
		UseCase:  useCase,
		Validate: validator.New(),
	}
}

// GetMethods handles get all payment methods
func (a *PaymentAdaptor) GetMethods(w http.ResponseWriter, r *http.Request) {
	methods, err := a.UseCase.GetPaymentMethods(r.Context())
	if err != nil {
		utils.ResponseInternalError(w, "failed to get payment methods")
		return
	}

	utils.ResponseOK(w, "success get payment methods", methods)
}

// ProcessPayment handles payment processing
func (a *PaymentAdaptor) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ResponseUnauthorized(w, "unauthorized")
		return
	}

	var req dto.PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := a.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, err.Error())
		return
	}

	payment, err := a.UseCase.ProcessPayment(r.Context(), userID, req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseOK(w, "payment successful", payment)
}
