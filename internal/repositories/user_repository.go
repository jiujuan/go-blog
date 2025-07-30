package repositories

import (
	"go-blog/internal/database"
	"go-blog/internal/models"
)

type userRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *userRepository) Create(user *models.User) error {
	return r.BaseRepository.Create(user)
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.BaseRepository.GetByID(&user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.GetDB().GetByField(&user, "email", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.GetDB().GetByField(&user, "username", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.BaseRepository.Update(user)
}

func (r *userRepository) Delete(id uint) error {
	return r.BaseRepository.Delete(&models.User{}, id)
}