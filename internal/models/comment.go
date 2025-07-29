package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ArticleID uint           `json:"article_id" gorm:"not null" validate:"required,min=1"`
	Article   Article        `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	UserID    uint           `json:"user_id" gorm:"not null" validate:"required,min=1"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Content   string         `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=2000"`
	ParentID  *uint          `json:"parent_id" validate:"omitempty,min=1"`
	Parent    *Comment       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies   []Comment      `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Comment model
func (Comment) TableName() string {
	return "comments"
}

// Validate validates the Comment model
func (c *Comment) Validate() error {
	if err := ValidateStruct(c); err != nil {
		return err
	}
	
	// Additional custom validations
	if err := c.validateContent(); err != nil {
		return err
	}
	
	if err := c.validateParentComment(); err != nil {
		return err
	}
	
	return nil
}

// validateContent performs additional content validation
func (c *Comment) validateContent() error {
	content := strings.TrimSpace(c.Content)
	if len(content) == 0 {
		return errors.New("content cannot be empty or contain only whitespace")
	}
	
	// Check for minimum meaningful content length
	if len(content) < 3 {
		return errors.New("comment content must be at least 3 characters long")
	}
	
	return nil
}

// validateParentComment validates parent comment relationship
func (c *Comment) validateParentComment() error {
	// If this is a reply, ensure it's not replying to itself
	if c.ParentID != nil && *c.ParentID == c.ID {
		return errors.New("comment cannot be a reply to itself")
	}
	
	return nil
}

// BeforeCreate hook for GORM
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	return c.Validate()
}

// BeforeUpdate hook for GORM
func (c *Comment) BeforeUpdate(tx *gorm.DB) error {
	return c.Validate()
}