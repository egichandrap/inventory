package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/example/jwt-ddd-clean/internal/handler"
)

// TokenHTTPHandler handles HTTP requests for token operations
type TokenHTTPHandler struct {
	tokenHandler *handler.TokenHandler
}

// NewTokenHTTPHandler creates a new TokenHTTPHandler
func NewTokenHTTPHandler(tokenHandler *handler.TokenHandler) *TokenHTTPHandler {
	return &TokenHTTPHandler{
		tokenHandler: tokenHandler,
	}
}

// GenerateTokenRequest represents the HTTP request body for token generation
type GenerateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshTokenRequest represents the HTTP request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RevokeTokenRequest represents the HTTP request body for token revocation
type RevokeTokenRequest struct {
	Token string `json:"token"`
}

// GenerateToken handles POST /api/token/generate
func (h *TokenHTTPHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Method not allowed"))
		return
	}

	var req GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Invalid request body"))
		return
	}

	if req.Username == "" || req.Password == "" {
		h.sendError(w, apperrors.ErrMissingFieldErr.WithDetails("Username and password are required"))
		return
	}

	response, err := h.tokenHandler.GenerateToken(req.Username, req.Password)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, "Token generated successfully", response, http.StatusOK)
}

// RefreshToken handles POST /api/token/refresh
func (h *TokenHTTPHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Method not allowed"))
		return
	}

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Invalid request body"))
		return
	}

	if req.RefreshToken == "" {
		h.sendError(w, apperrors.ErrMissingFieldErr.WithDetails("Refresh token is required"))
		return
	}

	response, err := h.tokenHandler.RefreshToken(req.RefreshToken)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, "Token refreshed successfully", response, http.StatusOK)
}

// ValidateToken handles POST /api/token/validate
func (h *TokenHTTPHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Method not allowed"))
		return
	}

	// Try to get token from header first
	token := r.Header.Get("Authorization")
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	// If no header, try request body
	if token == "" {
		var req struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.sendError(w, apperrors.ErrValidationErr.WithDetails("Invalid request body"))
			return
		}
		token = req.Token
	}

	if token == "" {
		h.sendError(w, apperrors.ErrMissingFieldErr.WithDetails("Token is required"))
		return
	}

	response, err := h.tokenHandler.ValidateToken(token)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, "Token validation result", response, http.StatusOK)
}

// RevokeToken handles POST /api/token/revoke
func (h *TokenHTTPHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Method not allowed"))
		return
	}

	var req RevokeTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Invalid request body"))
		return
	}

	if req.Token == "" {
		h.sendError(w, apperrors.ErrMissingFieldErr.WithDetails("Token is required"))
		return
	}

	if err := h.tokenHandler.RevokeToken(req.Token); err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, "Token revoked successfully", nil, http.StatusOK)
}

// Health handles GET /api/health
func (h *TokenHTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.sendSuccess(w, "Service is healthy", map[string]string{
		"status":  "up",
		"service": "jwt-ddd-clean",
	}, http.StatusOK)
}

func (h *TokenHTTPHandler) sendSuccess(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := apperrors.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *TokenHTTPHandler) sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		w.WriteHeader(appErr.GetHTTPStatus())
		json.NewEncoder(w).Encode(appErr.ToResponse())
		return
	}

	// Fallback for non-AppError
	w.WriteHeader(http.StatusInternalServerError)
	response := apperrors.ErrorResponse{
		Success: false,
		Error: apperrors.ErrorDetail{
			Code:    string(apperrors.ErrInternal),
			Message: "An unexpected error occurred",
			Details: err.Error(),
		},
	}
	json.NewEncoder(w).Encode(response)
}
