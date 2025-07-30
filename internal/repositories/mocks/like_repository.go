package mocks

import (
	"go-blog/internal/models"

	"github.com/stretchr/testify/mock"
)

// LikeRepository is a mock implementation of repositories.LikeRepository
type LikeRepository struct {
	mock.Mock
}

func (m *LikeRepository) Create(like *models.Like) error {
	args := m.Called(like)
	return args.Error(0)
}

func (m *LikeRepository) Delete(userID, articleID uint) error {
	args := m.Called(userID, articleID)
	return args.Error(0)
}

func (m *LikeRepository) GetByUserAndArticle(userID, articleID uint) (*models.Like, error) {
	args := m.Called(userID, articleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Like), args.Error(1)
}

func (m *LikeRepository) CountByArticle(articleID uint) (int64, error) {
	args := m.Called(articleID)
	return args.Get(0).(int64), args.Error(1)
}