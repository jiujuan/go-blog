package services

import (
	"errors"
	"go-blog/internal/models"
	"go-blog/internal/repositories"
	"gorm.io/gorm"
)

type LikeService struct {
	likeRepo    repositories.LikeRepository
	articleRepo repositories.ArticleRepository
	userRepo    repositories.UserRepository
}

// NewLikeService creates a new like service
func NewLikeService(
	likeRepo repositories.LikeRepository,
	articleRepo repositories.ArticleRepository,
	userRepo repositories.UserRepository,
) *LikeService {
	return &LikeService{
		likeRepo:    likeRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

// ToggleLike toggles like status for an article by a user
func (s *LikeService) ToggleLike(userID, articleID uint) (bool, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("user not found")
		}
		return false, err
	}

	// Verify article exists
	_, err = s.articleRepo.GetByID(articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("article not found")
		}
		return false, err
	}

	// Check if like already exists
	existingLike, err := s.likeRepo.GetByUserAndArticle(userID, articleID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if existingLike != nil {
		// Unlike - remove the like
		err = s.likeRepo.Delete(userID, articleID)
		if err != nil {
			return false, err
		}
		return false, nil // false means unliked
	} else {
		// Like - create new like
		like := &models.Like{
			UserID:    userID,
			ArticleID: articleID,
		}
		err = s.likeRepo.Create(like)
		if err != nil {
			return false, err
		}
		return true, nil // true means liked
	}
}

// IsLikedByUser checks if an article is liked by a specific user
func (s *LikeService) IsLikedByUser(userID, articleID uint) (bool, error) {
	like, err := s.likeRepo.GetByUserAndArticle(userID, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return like != nil, nil
}

// GetLikeCount returns the total number of likes for an article
func (s *LikeService) GetLikeCount(articleID uint) (int64, error) {
	return s.likeRepo.CountByArticle(articleID)
}

// GetArticleLikeStatus returns like count and whether user has liked the article
func (s *LikeService) GetArticleLikeStatus(articleID uint, userID *uint) (int64, bool, error) {
	count, err := s.GetLikeCount(articleID)
	if err != nil {
		return 0, false, err
	}

	var isLiked bool
	if userID != nil {
		isLiked, err = s.IsLikedByUser(*userID, articleID)
		if err != nil {
			return count, false, err
		}
	}

	return count, isLiked, nil
}