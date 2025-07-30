package mocks

import (
	"go-blog/internal/models"

	"github.com/stretchr/testify/mock"
)

// ArticleRepository is a mock implementation of repositories.ArticleRepository
type ArticleRepository struct {
	mock.Mock
}

func (m *ArticleRepository) Create(article *models.Article) error {
	args := m.Called(article)
	return args.Error(0)
}

func (m *ArticleRepository) GetByID(id uint) (*models.Article, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *ArticleRepository) GetBySlug(slug string) (*models.Article, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *ArticleRepository) List(offset, limit int, filters map[string]interface{}) ([]models.Article, int64, error) {
	args := m.Called(offset, limit, filters)
	return args.Get(0).([]models.Article), args.Get(1).(int64), args.Error(2)
}

func (m *ArticleRepository) Update(article *models.Article) error {
	args := m.Called(article)
	return args.Error(0)
}

func (m *ArticleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ArticleRepository) Search(query string, offset, limit int) ([]models.Article, int64, error) {
	args := m.Called(query, offset, limit)
	return args.Get(0).([]models.Article), args.Get(1).(int64), args.Error(2)
}

func (m *ArticleRepository) GetArchive() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *ArticleRepository) GetByMonth(year, month int, offset, limit int) ([]models.Article, int64, error) {
	args := m.Called(year, month, offset, limit)
	return args.Get(0).([]models.Article), args.Get(1).(int64), args.Error(2)
}

func (m *ArticleRepository) GetByAuthorID(authorID uint, limit, offset int) ([]*models.Article, error) {
	args := m.Called(authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Article), args.Error(1)
}

func (m *ArticleRepository) CountByAuthorID(authorID uint) (int64, error) {
	args := m.Called(authorID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *ArticleRepository) IncrementViewCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ArticleRepository) UpdateStatistics(id uint, viewCount, likeCount, commentCount uint) error {
	args := m.Called(id, viewCount, likeCount, commentCount)
	return args.Error(0)
}

func (m *ArticleRepository) GetStatistics(id uint) (viewCount, likeCount, commentCount uint, err error) {
	args := m.Called(id)
	return args.Get(0).(uint), args.Get(1).(uint), args.Get(2).(uint), args.Error(3)
}