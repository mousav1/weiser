package repositories

import (
	"errors"

	"github.com/mousav1/weiser/app/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db             *gorm.DB
	BaseRepository BaseRepository[models.User]
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.User](db),
	}
}

// GetByUsername retrieves a user by its username.
func (ur *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user *models.User
	err := ur.BaseRepository.GetFirst(user, "username = ?", username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieves a user by its email.
func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user *models.User
	err := ur.BaseRepository.GetFirst(user, "email = ?", email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// func (br *UserRepository) Create(record *models.User) error {
// 	return br.baseRepository.Create(record)
// }

// // Update updates an existing record in the database.
// func (br *UserRepository) Update(record *models.User) error {
// 	return br.baseRepository.Update(record)
// }

// func (br *UserRepository) Delete(record *models.User) error {
// 	return br.baseRepository.Delete(record)
// }

// func (ur *UserRepository) GetByID(id uint) (*models.User, error) {
// 	user := &models.User{}
// 	user, err := ur.baseRepository.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }
