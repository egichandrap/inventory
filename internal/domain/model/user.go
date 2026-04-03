package model

import "time"

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
	RoleAdmin      UserRole = "ADMIN"
	RoleCashier    UserRole = "CASHIER"
	RoleViewer     UserRole = "VIEWER"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
	StatusSuspended UserStatus = "SUSPENDED"
)

// User represents a user entity in the domain
type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Password     string     `json:"-"` // Never serialize password
	PasswordHash string     `json:"-"` // Hashed password
	Email        string     `json:"email"`
	FullName     string     `json:"full_name"`
	Role         UserRole   `json:"role"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// NewUser creates a new user entity
func NewUser(username, email, fullName string, role UserRole) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Email:     email,
		FullName:  fullName,
		Role:      role,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CanManageInventory checks if user has permission to manage inventory
func (u *User) CanManageInventory() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin
}

// CanAccessPOS checks if user can access POS features
func (u *User) CanAccessPOS() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAdmin || u.Role == RoleCashier
}

// CanManageUsers checks if user can manage other users
func (u *User) CanManageUsers() bool {
	return u.Role == RoleSuperAdmin
}

// IsActive checks if user account is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}
