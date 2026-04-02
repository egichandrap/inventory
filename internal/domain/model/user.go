package model

// User represents a user entity in the domain
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Never serialize password
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// NewUser creates a new user entity
func NewUser(id, username, email, role string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
	}
}
