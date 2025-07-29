package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusPublished ArticleStatus = "published"
	StatusArchived  ArticleStatus = "archived"
)

type Article struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"size:255;not null" validate:"required,min=1,max=255"`
	Slug         string         `json:"slug" gorm:"uniqueIndex;size:255;not null" validate:"required,slug,max=255"`
	Content      string         `json:"content" gorm:"type:longtext;not null" validate:"required,min=1"`
	Excerpt      string         `json:"excerpt" gorm:"type:text" validate:"omitempty,max=500"`
	AuthorID     uint           `json:"author_id" gorm:"not null" validate:"required,min=1"`
	Author       User           `json:"author" gorm:"foreignKey:AuthorID"`
	CategoryID   *uint          `json:"category_id" validate:"omitempty,min=1"`
	Category     *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags         []Tag          `json:"tags,omitempty" gorm:"many2many:article_tags"`
	Comments     []Comment      `json:"comments,omitempty"`
	Likes        []Like         `json:"likes,omitempty"`
	Status       ArticleStatus  `json:"status" gorm:"type:enum('draft','published','archived');default:'draft'" validate:"required,article_status"`
	ViewCount    uint           `json:"view_count" gorm:"default:0"`
	LikeCount    uint           `json:"like_count" gorm:"default:0"`
	CommentCount uint           `json:"comment_count" gorm:"default:0"`
	PublishedAt  *time.Time     `json:"published_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Article model
func (Article) TableName() string {
	return "articles"
}

// Validate validates the Article model
func (a *Article) Validate() error {
	if err := ValidateStruct(a); err != nil {
		return err
	}
	
	// Additional custom validations
	if err := a.validateTitle(); err != nil {
		return err
	}
	
	if err := a.validateContent(); err != nil {
		return err
	}
	
	if err := a.validatePublishedStatus(); err != nil {
		return err
	}
	
	return nil
}

// validateTitle performs additional title validation
func (a *Article) validateTitle() error {
	title := strings.TrimSpace(a.Title)
	if title != a.Title {
		return errors.New("title cannot have leading or trailing spaces")
	}
	
	if strings.Contains(title, "\n") || strings.Contains(title, "\r") {
		return errors.New("title cannot contain line breaks")
	}
	
	return nil
}

// validateContent performs additional content validation
func (a *Article) validateContent() error {
	content := strings.TrimSpace(a.Content)
	if len(content) == 0 {
		return errors.New("content cannot be empty or contain only whitespace")
	}
	
	return nil
}

// validatePublishedStatus validates published status requirements
func (a *Article) validatePublishedStatus() error {
	if a.Status == StatusPublished {
		if a.PublishedAt == nil {
			now := time.Now()
			a.PublishedAt = &now
		}
	}
	
	return nil
}

// BeforeCreate hook for GORM
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	return a.Validate()
}

// BeforeUpdate hook for GORM
func (a *Article) BeforeUpdate(tx *gorm.DB) error {
	return a.Validate()
}