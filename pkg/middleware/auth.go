package middleware

import (
	"context"
	"net/http"
	"project-app-bioskop/internal/usecase"
	"project-app-bioskop/pkg/utils"
)

// AuthMiddleware provides token-based authentication
type AuthMiddleware struct {
	AuthUseCase usecase.AuthUseCaseInterface
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(authUseCase usecase.AuthUseCaseInterface) *AuthMiddleware {
	return &AuthMiddleware{AuthUseCase: authUseCase}
}

// RequireAuth validates token and injects user ID into context
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.ResponseBadRequest(w, http.StatusUnauthorized, "missing authorization token", nil)
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		user, err := m.AuthUseCase.ValidateToken(r.Context(), token)
		if err != nil {
			utils.ResponseBadRequest(w, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Auth is a legacy middleware for backward compatibility
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ID := "2"
		ctx = context.WithValue(ctx, "ctxid", ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
