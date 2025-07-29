package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null" validate:"required,min=1"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	ArticleID uint           `json:"article_id" gorm:"not null" validate:"required,min=1"`
	Article   Article        `json:"article" gorm:"foreignKey:ArticleID"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Like model
func (Like) TableName() string {
	return "likes"
}

// Validate validates the Like model
func (l *Like) Validate() error {
	if err := ValidateStruct(l); err != nil {
		return err
	}
	
	// Additional custom validations
	if err := l.validateUserArticleCombination(); err != nil {
		return err
	}
	
	return nil
}

// validateUserArticleCombination validates the user-article relationship
func (l *Like) validateUserArticleCombination() error {
	if l.UserID == 0 {
		return errors.New("user ID is required")
	}
	
	if l.ArticleID == 0 {
		return errors.New("article ID is required")
	}
	
	return nil
}

// BeforeCreate hook for GORM
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	if err := l.Validate(); err != nil {
		return err
	}
	
	// Check for existing like to prevent duplicates
	var existingLike Like
	result := tx.Where("user_id = ? AND article_id = ?", l.UserID, l.ArticleID).First(&existingLike)
	if result.Error == nil {
		return errors.New("user has already liked this article")
	}
	
	return nil
}

// BeforeUpdate hook for GORM
func (l *Like) BeforeUpdate(tx *gorm.DB) error {
	return l.Validate()
}