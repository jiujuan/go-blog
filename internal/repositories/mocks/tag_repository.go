package mocks

import (
	"go-blog/internal/models"

	"github.com/stretchr/testify/mock"
)

// TagRepository is a mock implementation of repositories.TagRepository
type TagRepository struct {
	mock.Mock
}

func (m *TagRepository) Create(tag *models.Tag) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *TagRepository) GetByID(id uint) (*models.Tag, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *TagRepository) GetBySlug(slug string) (*models.Tag, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *TagRepository) GetByName(name string) (*models.Tag, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *TagRepository) List() ([]models.Tag, error) {
	args := m.Called()
	return args.Get(0).([]models.Tag), args.Error(1)
}

func (m *TagRepository) GetArticles(tagID uint, offset, limit int) ([]models.Article, int64, error) {
	args := m.Called(tagID, offset, limit)
	return args.Get(0).([]models.Article), args.Get(1).(int64), args.Error(2)
}