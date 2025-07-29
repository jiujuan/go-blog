package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:100;not null" validate:"required,min=1,max=100"`
	Description string         `json:"description" gorm:"type:text" validate:"omitempty,max=1000"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;size:100;not null" validate:"required,slug,max=100"`
	Articles    []Article      `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Category model
func (Category) TableName() string {
	return "categories"
}

// Validate validates the Category model
func (c *Category) Validate() error {
	if err := ValidateStruct(c); err != nil {
		return err
	}
	
	// Additional custom validations
	if err := c.validateName(); err != nil {
		return err
	}
	
	return nil
}

// validateName performs additional name validation
func (c *Category) validateName() error {
	name := strings.TrimSpace(c.Name)
	if name != c.Name {
		return errors.New("name cannot have leading or trailing spaces")
	}
	
	if strings.Contains(name, "\n") || strings.Contains(name, "\r") {
		return errors.New("name cannot contain line breaks")
	}
	
	// Check for reserved category names
	reservedNames := []string{"admin", "api", "www", "blog", "category", "categories"}
	for _, reserved := range reservedNames {
		if strings.ToLower(name) == reserved {
			return errors.New("category name is reserved and cannot be used")
		}
	}
	
	return nil
}

// BeforeCreate hook for GORM
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	return c.Validate()
}

// BeforeUpdate hook for GORM
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	return c.Validate()
}