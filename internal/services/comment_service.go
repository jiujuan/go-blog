package services

import (
	"errors"
	"go-blog/internal/models"
	"go-blog/internal/repositories"
	"gorm.io/gorm"
)

type CommentService struct {
	commentRepo repositories.CommentRepository
	articleRepo repositories.ArticleRepository
	userRepo    repositories.UserRepository
}

// NewCommentService creates a new comment service
func NewCommentService(
	commentRepo repositories.CommentRepository,
	articleRepo repositories.ArticleRepository,
	userRepo repositories.UserRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new comment with validation
func (s *CommentService) Create(comment *models.Comment) error {
	// Verify user exists
	_, err := s.userRepo.GetByID(comment.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Verify article exists
	_, err = s.articleRepo.GetByID(comment.ArticleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("article not found")
		}
		return err
	}

	// If this is a reply, verify parent comment exists and belongs to same article
	if comment.ParentID != nil {
		parentComment, err := s.commentRepo.GetByID(*comment.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parent comment not found")
			}
			return err
		}

		// Ensure parent comment belongs to the same article
		if parentComment.ArticleID != comment.ArticleID {
			return errors.New("parent comment must belong to the same article")
		}
	}

	return s.commentRepo.Create(comment)
}

// GetByArticle retrieves comments for an article with threading
func (s *CommentService) GetByArticle(articleID uint) ([]models.Comment, error) {
	// Verify article exists
	_, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}

	return s.commentRepo.GetByArticle(articleID)
}

// GetByID retrieves a comment by ID
func (s *CommentService) GetByID(id uint) (*models.Comment, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	return comment, nil
}

// Update updates a comment with authorization check
func (s *CommentService) Update(commentID uint, userID uint, content string) (*models.Comment, error) {
	// Get existing comment
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	// Check if user is the author of the comment
	if comment.UserID != userID {
		return nil, errors.New("unauthorized: can only update your own comments")
	}

	// Update content
	comment.Content = content
	err = s.commentRepo.Update(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// Delete deletes a comment with authorization check
func (s *CommentService) Delete(commentID uint, userID uint) error {
	// Get existing comment
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("comment not found")
		}
		return err
	}

	// Check if user is the author of the comment
	if comment.UserID != userID {
		return errors.New("unauthorized: can only delete your own comments")
	}

	return s.commentRepo.Delete(commentID)
}