package repositories

import (
	"go-blog/internal/database"
	"go-blog/internal/models"
)

type commentRepository struct {
	*BaseRepository
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *database.DB) CommentRepository {
	return &commentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.BaseRepository.Create(comment)
}

func (r *commentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.BaseRepository.GetByID(&comment, id, "User", "Replies")
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) GetByArticle(articleID uint) ([]models.Comment, error) {
	var comments []models.Comment
	
	// For complex queries with conditions, use the underlying GORM DB
	err := r.GetDB().GetDB().Preload("User").Preload("Replies").
		Where("article_id = ? AND parent_id IS NULL", articleID).
		Order("created_at ASC").Find(&comments).Error
	return comments, err
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.BaseRepository.Update(comment)
}

func (r *commentRepository) Delete(id uint) error {
	return r.BaseRepository.Delete(&models.Comment{}, id)
}