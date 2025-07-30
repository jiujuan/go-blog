package mocks

import (
	"go-blog/internal/models"

	"github.com/stretchr/testify/mock"
)

// CategoryRepository is a mock implementation of repositories.CategoryRepository
type CategoryRepository struct {
	mock.Mock
}

func (m *CategoryRepository) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryRepository) GetByID(id uint) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepository) List() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *CategoryRepository) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *CategoryRepository) GetArticles(categoryID uint, offset, limit int) ([]models.Article, int64, error) {
	args := m.Called(categoryID, offset, limit)
	return args.Get(0).([]models.Article), args.Get(1).(int64), args.Error(2)
}