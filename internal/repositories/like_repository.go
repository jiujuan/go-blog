package repositories

import (
	"go-blog/internal/database"
	"go-blog/internal/models"
)

type likeRepository struct {
	*BaseRepository
}

// NewLikeRepository creates a new like repository
func NewLikeRepository(db *database.DB) LikeRepository {
	return &likeRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *likeRepository) Create(like *models.Like) error {
	return r.BaseRepository.Create(like)
}

func (r *likeRepository) Delete(userID, articleID uint) error {
	return r.GetDB().BulkDelete(&models.Like{}, "user_id = ? AND article_id = ?", userID, articleID)
}

func (r *likeRepository) GetByUserAndArticle(userID, articleID uint) (*models.Like, error) {
	var like models.Like
	
	// Use raw GORM for complex WHERE conditions
	err := r.GetDB().GetDB().Where("user_id = ? AND article_id = ?", userID, articleID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) CountByArticle(articleID uint) (int64, error) {
	return r.BaseRepository.Count(&models.Like{}, "article_id = ?", articleID)
}