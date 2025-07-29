package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,min=1,max=50"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;size:50;not null" validate:"required,slug,max=50"`
	Articles  []Article      `json:"articles,omitempty" gorm:"many2many:article_tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Tag model
func (Tag) TableName() string {
	return "tags"
}

// Validate validates the Tag model
func (t *Tag) Validate() error {
	if err := ValidateStruct(t); err != nil {
		return err
	}
	
	// Additional custom validations
	if err := t.validateName(); err != nil {
		return err
	}
	
	return nil
}

// validateName performs additional name validation
func (t *Tag) validateName() error {
	name := strings.TrimSpace(t.Name)
	if name != t.Name {
		return errors.New("name cannot have leading or trailing spaces")
	}
	
	if strings.Contains(name, "\n") || strings.Contains(name, "\r") {
		return errors.New("name cannot contain line breaks")
	}
	
	// Tags should not contain special characters except hyphens and underscores
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || char == '-' || char == '_' || char == ' ') {
			return errors.New("tag name can only contain letters, numbers, spaces, hyphens, and underscores")
		}
	}
	
	return nil
}

// BeforeCreate hook for GORM
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	return t.Validate()
}

// BeforeUpdate hook for GORM
func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
	return t.Validate()
}