package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Category represents a product category aggregate root
type Category struct {
	id          string
	name        string
	slug        string
	description string
	parentID    string // for hierarchical categories
	level       int
	sortOrder   int
	isActive    bool
	productCount int
	createdAt   time.Time
	updatedAt   time.Time
}

// NewCategory creates a new product category
func NewCategory(name, slug, description string, parentID string) (*Category, error) {
	if name == "" {
		return nil, fmt.Errorf("nama kategori tidak boleh kosong")
	}
	if slug == "" {
		return nil, fmt.Errorf("slug kategori tidak boleh kosong")
	}

	now := time.Now()
	level := 0
	if parentID != "" {
		level = 1 // Can be enhanced to support deeper hierarchies
	}

	return &Category{
		id:           uuid.New().String(),
		name:         name,
		slug:         slug,
		description:  description,
		parentID:     parentID,
		level:        level,
		sortOrder:    0,
		isActive:     true,
		productCount: 0,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructCategory recreates a category from database
func ReconstructCategory(
	id, name, slug, description, parentID string,
	level, sortOrder, productCount int,
	isActive bool,
	createdAt, updatedAt time.Time,
) *Category {
	return &Category{
		id:           id,
		name:         name,
		slug:         slug,
		description:  description,
		parentID:     parentID,
		level:        level,
		sortOrder:    sortOrder,
		isActive:     isActive,
		productCount: productCount,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Accessors
func (c *Category) ID() string            { return c.id }
func (c *Category) Name() string          { return c.name }
func (c *Category) Slug() string          { return c.slug }
func (c *Category) Description() string   { return c.description }
func (c *Category) ParentID() string      { return c.parentID }
func (c *Category) Level() int            { return c.level }
func (c *Category) SortOrder() int        { return c.sortOrder }
func (c *Category) IsActive() bool        { return c.isActive }
func (c *Category) ProductCount() int     { return c.productCount }
func (c *Category) CreatedAt() time.Time  { return c.createdAt }
func (c *Category) UpdatedAt() time.Time  { return c.updatedAt }
func (c *Category) IsRoot() bool          { return c.parentID == "" }

// UpdateDetails updates category details
func (c *Category) UpdateDetails(name, slug, description string) error {
	if name == "" {
		return fmt.Errorf("nama kategori tidak boleh kosong")
	}
	if slug == "" {
		return fmt.Errorf("slug kategori tidak boleh kosong")
	}

	c.name = name
	c.slug = slug
	c.description = description
	c.updatedAt = time.Now()
	return nil
}

// Activate activates the category
func (c *Category) Activate() {
	c.isActive = true
	c.updatedAt = time.Now()
}

// Deactivate deactivates the category
func (c *Category) Deactivate() {
	c.isActive = false
	c.updatedAt = time.Now()
}

// IncrementProductCount increments the product count
func (c *Category) IncrementProductCount() {
	c.productCount++
	c.updatedAt = time.Now()
}

// DecrementProductCount decrements the product count
func (c *Category) DecrementProductCount() error {
	if c.productCount <= 0 {
		return fmt.Errorf("product count tidak boleh negatif")
	}
	c.productCount--
	c.updatedAt = time.Now()
	return nil
}

// UpdateSortOrder updates the sort order
func (c *Category) UpdateSortOrder(order int) {
	c.sortOrder = order
	c.updatedAt = time.Now()
}

// IsChildOf checks if this category is a child of another
func (c *Category) IsChildOf(parentID string) bool {
	return c.parentID == parentID
}
