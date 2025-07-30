package repositories

import (
	"go-blog/internal/database"
	"go-blog/internal/models"
)

type categoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *database.DB) CategoryRepository {
	return &categoryRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.BaseRepository.Create(category)
}

func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.BaseRepository.GetByID(&category, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.GetDB().GetByField(&category, "slug", slug)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	options := &database.QueryOptions{
		Page:    1,
		Limit:   1000, // Large limit for getting all categories
		OrderBy: "name ASC",
	}
	
	_, err := r.BaseRepository.List(&categories, options)
	return categories, err
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.BaseRepository.Update(category)
}

func (r *categoryRepository) Delete(id uint) error {
	return r.BaseRepository.Delete(&models.Category{}, id)
}

func (r *categoryRepository) GetArticles(categoryID uint, offset, limit int) ([]models.Article, int64, error) {
	var articles []models.Article
	
	// Convert offset/limit to page-based pagination
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	
	options := &database.QueryOptions{
		Page:    page,
		Limit:   limit,
		OrderBy: "created_at DESC",
		Filters: map[string]interface{}{
			"category_id": categoryID,
		},
		Preloads: []string{"Author", "Category", "Tags"},
	}
	
	result, err := r.BaseRepository.List(&articles, options)
	if err != nil {
		return nil, 0, err
	}
	
	return articles, result.Total, nil
}