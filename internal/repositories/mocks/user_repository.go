package mocks

import (
	"go-blog/internal/models"

	"github.com/stretchr/testify/mock"
)

// UserRepository is a mock implementation of repositories.UserRepository
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}