package repositories

import (
	"go-blog/internal/database"
	"go-blog/internal/models"
)

type tagRepository struct {
	*BaseRepository
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *database.DB) TagRepository {
	return &tagRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.BaseRepository.Create(tag)
}

func (r *tagRepository) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.BaseRepository.GetByID(&tag, id)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.GetDB().GetByField(&tag, "slug", slug)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.GetDB().GetByField(&tag, "name", name)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) List() ([]models.Tag, error) {
	var tags []models.Tag
	options := &database.QueryOptions{
		Page:    1,
		Limit:   1000, // Large limit for getting all tags
		OrderBy: "name ASC",
	}
	
	_, err := r.BaseRepository.List(&tags, options)
	return tags, err
}

func (r *tagRepository) GetArticles(tagID uint, offset, limit int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64
	
	// For many-to-many relationships, we need to use raw GORM query
	// Convert offset/limit to page calculation
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	
	query := r.GetDB().GetDB().Model(&models.Article{}).
		Preload("Author").Preload("Category").Preload("Tags").
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ?", tagID)

	// Count total records
	query.Count(&total)
	
	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&articles).Error
	return articles, total, err
}