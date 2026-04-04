package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Store represents a store/branch aggregate root
type Store struct {
	id          string
	code        string
	name        string
	address     string
	phone       string
	email       string
	managerID   string
	managerName string
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

// NewStore creates a new store
func NewStore(code, name, address, phone, email string) (*Store, error) {
	if code == "" {
		return nil, fmt.Errorf("kode toko tidak boleh kosong")
	}
	if name == "" {
		return nil, fmt.Errorf("nama toko tidak boleh kosong")
	}

	now := time.Now()
	return &Store{
		id:         uuid.New().String(),
		code:       code,
		name:       name,
		address:    address,
		phone:      phone,
		email:      email,
		isActive:   true,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// ReconstructStore recreates a store from database
func ReconstructStore(
	id, code, name, address, phone, email, managerID, managerName string,
	isActive bool,
	createdAt, updatedAt time.Time,
) *Store {
	return &Store{
		id:          id,
		code:        code,
		name:        name,
		address:     address,
		phone:       phone,
		email:       email,
		managerID:   managerID,
		managerName: managerName,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Accessors
func (s *Store) ID() string           { return s.id }
func (s *Store) Code() string         { return s.code }
func (s *Store) Name() string         { return s.name }
func (s *Store) Address() string      { return s.address }
func (s *Store) Phone() string        { return s.phone }
func (s *Store) Email() string        { return s.email }
func (s *Store) ManagerID() string    { return s.managerID }
func (s *Store) ManagerName() string  { return s.managerName }
func (s *Store) IsActive() bool       { return s.isActive }
func (s *Store) CreatedAt() time.Time { return s.createdAt }
func (s *Store) UpdatedAt() time.Time { return s.updatedAt }

// UpdateDetails updates store details
func (s *Store) UpdateDetails(name, address, phone, email string) {
	s.name = name
	s.address = address
	s.phone = phone
	s.email = email
	s.updatedAt = time.Now()
}

// AssignManager assigns a manager to the store
func (s *Store) AssignManager(managerID, managerName string) {
	s.managerID = managerID
	s.managerName = managerName
	s.updatedAt = time.Now()
}

// Activate activates the store
func (s *Store) Activate() {
	s.isActive = true
	s.updatedAt = time.Now()
}

// Deactivate deactivates the store
func (s *Store) Deactivate() {
	s.isActive = false
	s.updatedAt = time.Now()
}
