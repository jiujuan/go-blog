package mocks

import (
	"go-blog/internal/models"

	"github.com/stretchr/testify/mock"
)

// CommentRepository is a mock implementation of repositories.CommentRepository
type CommentRepository struct {
	mock.Mock
}

func (m *CommentRepository) Create(comment *models.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *CommentRepository) GetByID(id uint) (*models.Comment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *CommentRepository) GetByArticle(articleID uint) ([]models.Comment, error) {
	args := m.Called(articleID)
	return args.Get(0).([]models.Comment), args.Error(1)
}

func (m *CommentRepository) Update(comment *models.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *CommentRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}