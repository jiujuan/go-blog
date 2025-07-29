package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,username"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email,max=100"`
	Password  string         `json:"-" gorm:"size:255;not null;column:password_hash" validate:"required,min=8,max=255"`
	AvatarURL string         `json:"avatar_url" gorm:"size:255;column:avatar_url" validate:"omitempty,url,max=255"`
	Bio       string         `json:"bio" gorm:"type:text" validate:"omitempty,max=1000"`
	Articles  []Article      `json:"articles,omitempty" gorm:"foreignKey:AuthorID"`
	Comments  []Comment      `json:"comments,omitempty"`
	Likes     []Like         `json:"likes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// Validate validates the User model
func (u *User) Validate() error {
	if err := ValidateStruct(u); err != nil {
		return err
	}
	
	// Additional custom validations
	if err := u.validateUsername(); err != nil {
		return err
	}
	
	return nil
}

// validateUsername performs additional username validation
func (u *User) validateUsername() error {
	username := strings.TrimSpace(u.Username)
	if username != u.Username {
		return errors.New("username cannot have leading or trailing spaces")
	}
	
	// Check for reserved usernames
	reservedUsernames := []string{"admin", "root", "api", "www", "mail", "ftp", "blog", "user", "test"}
	for _, reserved := range reservedUsernames {
		if strings.ToLower(username) == reserved {
			return errors.New("username is reserved and cannot be used")
		}
	}
	
	return nil
}

// BeforeCreate hook for GORM
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.Validate()
}

// BeforeUpdate hook for GORM
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}