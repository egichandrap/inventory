package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	"github.com/example/jwt-ddd-clean/internal/dto"
	"github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	jwtProvider JWTProvider
}

// NewAuthService creates a new AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtProvider JWTProvider,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtProvider: jwtProvider,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewValidationError("username atau password salah")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.NewValidationError("akun tidak aktif atau telah ditangguhkan")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.NewValidationError("username atau password salah")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("failed to update last login: %v\n", err)
	}

	// Generate tokens
	accessToken, err := s.jwtProvider.GenerateTokenWithDuration(user.ID, user.Username, user.Role, 24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError("gagal membuat access token")
	}

	refreshToken, err := s.jwtProvider.GenerateTokenWithDuration(user.ID, user.Username, user.Role, 7*24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError("gagal membuat refresh token")
	}

	// Update user's last login
	now := time.Now()
	user.LastLoginAt = &now

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(24 * time.Hour.Seconds()),
		User:         dto.ToUserResponse(user),
	}, nil
}

// Logout invalidates user tokens
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	// Blacklist the access token (we'll use a simple approach with token string)
	// In production, you'd decode the token to get its expiration
	expiresAt := time.Now().Add(24 * time.Hour) // Default expiration
	if err := s.tokenRepo.Blacklist(ctx, accessToken, expiresAt); err != nil {
		return errors.NewInternalError("gagal melakukan logout")
	}

	return nil
}

// Register creates a new user
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if username exists
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewInternalError("gagal memeriksa username")
	}
	if exists {
		return nil, errors.NewValidationError("username telah digunakan")
	}

	// Check if email exists
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewInternalError("gagal memeriksa email")
	}
	if exists {
		return nil, errors.NewValidationError("email telah digunakan")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError("gagal memproses password")
	}

	// Create user
	user := model.NewUser(req.Username, req.Email, req.FullName, req.Role)
	user.PasswordHash = string(hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewInternalError("gagal membuat user: %v", err)
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtProvider.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.NewValidationError("refresh token tidak valid")
	}

	// Check if refresh token is blacklisted
	isBlacklisted, err := s.tokenRepo.IsBlacklisted(ctx, refreshToken)
	if err != nil {
		return nil, errors.NewInternalError("gagal memvalidasi refresh token")
	}
	if isBlacklisted {
		return nil, errors.NewValidationError("refresh token telah dicabut")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("user", "id", claims.UserID)
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.NewValidationError("akun tidak aktif atau telah ditangguhkan")
	}

	// Generate new access token
	accessToken, err := s.jwtProvider.GenerateTokenWithDuration(user.ID, user.Username, user.Role, 24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError("gagal membuat access token")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(24 * time.Hour.Seconds()),
		User:         dto.ToUserResponse(user),
	}, nil
}

// GetMe returns current user information
func (s *AuthService) GetMe(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user", "id", userID)
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewNotFoundError("user", "id", userID)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.NewValidationError("password lama salah")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalError("gagal memproses password")
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return errors.NewInternalError("gagal mengubah password")
	}

	return nil
}

// ListUsers returns paginated list of users
func (s *AuthService) ListUsers(ctx context.Context, filter repository.UserFilter) (*dto.UserListResponse, error) {
	paginatedUsers, err := s.userRepo.ListWithPagination(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil daftar user")
	}

	users := make([]dto.UserResponse, len(paginatedUsers.Users))
	for i, user := range paginatedUsers.Users {
		users[i] = dto.ToUserResponse(user)
	}

	return &dto.UserListResponse{
		Users:      users,
		Total:      paginatedUsers.Total,
		Limit:      paginatedUsers.Limit,
		Offset:     paginatedUsers.Offset,
		TotalPages: paginatedUsers.TotalPages,
	}, nil
}

// UpdateUser updates user information
func (s *AuthService) UpdateUser(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user", "id", userID)
	}

	// Update fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.NewInternalError("gagal mengupdate user")
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(ctx context.Context, userID string) error {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewNotFoundError("user", "id", userID)
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return errors.NewInternalError("gagal menghapus user")
	}

	return nil
}
