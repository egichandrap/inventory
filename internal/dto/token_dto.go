package dto

// TokenRequest represents a token generation request
type TokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RevokeTokenRequest represents a token revocation request
type RevokeTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// ValidateTokenResponse represents a token validation response
type ValidateTokenResponse struct {
	Valid    bool   `json:"valid"`
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}
