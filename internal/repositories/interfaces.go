package repositories

import (
	"time"

	"go-blog/internal/database"
	"go-blog/internal/models"
)

// UserRepository interface defines user data access methods
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

// SearchFilters represents advanced search filters
type SearchFilters struct {
	Status     string    `json:"status,omitempty"`
	CategoryID uint      `json:"category_id,omitempty"`
	AuthorID   uint      `json:"author_id,omitempty"`
	TagID      uint      `json:"tag_id,omitempty"`
	DateFrom   time.Time `json:"date_from,omitempty"`
	DateTo     time.Time `json:"date_to,omitempty"`
}

// ArticleRepository interface defines article data access methods
type ArticleRepository interface {
	Create(article *models.Article) error
	GetByID(id uint) (*models.Article, error)
	GetBySlug(slug string) (*models.Article, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.Article, int64, error)
	Update(article *models.Article) error
	Delete(id uint) error
	Search(query string, offset, limit int) ([]models.Article, int64, error)
	AdvancedSearch(query string, offset, limit int, filters *SearchFilters) ([]models.Article, int64, error)
	SearchWithBoolean(query string, offset, limit int, filters *SearchFilters) ([]models.Article, int64, error)
	GetArchive() (map[string]interface{}, error)
	GetByMonth(year, month int, offset, limit int) ([]models.Article, int64, error)
	GetByAuthorID(authorID uint, limit, offset int) ([]*models.Article, error)
	CountByAuthorID(authorID uint) (int64, error)
	IncrementViewCount(id uint) error
	UpdateStatistics(id uint, viewCount, likeCount, commentCount uint) error
	GetStatistics(id uint) (viewCount, likeCount, commentCount uint, err error)
}

// CategoryRepository interface defines category data access methods
type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	GetBySlug(slug string) (*models.Category, error)
	List() ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	GetArticles(categoryID uint, offset, limit int) ([]models.Article, int64, error)
}

// TagRepository interface defines tag data access methods
type TagRepository interface {
	Create(tag *models.Tag) error
	GetByID(id uint) (*models.Tag, error)
	GetBySlug(slug string) (*models.Tag, error)
	GetByName(name string) (*models.Tag, error)
	List() ([]models.Tag, error)
	GetArticles(tagID uint, offset, limit int) ([]models.Article, int64, error)
}

// CommentRepository interface defines comment data access methods
type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByID(id uint) (*models.Comment, error)
	GetByArticle(articleID uint) ([]models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uint) error
}

// LikeRepository interface defines like data access methods
type LikeRepository interface {
	Create(like *models.Like) error
	Delete(userID, articleID uint) error
	GetByUserAndArticle(userID, articleID uint) (*models.Like, error)
	CountByArticle(articleID uint) (int64, error)
}