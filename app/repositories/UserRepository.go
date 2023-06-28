package repositories

import (
	"errors"

	"github.com/mousav1/weiser/app/models"
	"gorm.io/gorm"
)

// UserRepository represents the repository for managing users.
type UserRepository interface {
	Create(*models.User) error
	GetByID(uint) (*models.User, error)
	GetByUsername(string) (*models.User, error)
	GetByEmail(string) (*models.User, error)
	Update(*models.User) error
	Delete(*models.User) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user.
func (ur *userRepository) Create(user *models.User) error {
	return ur.db.Create(user).Error
}

// GetByID retrieves a user by its ID.
func (ur *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := ur.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by its username.
func (ur *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by its email.
func (ur *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user.
func (ur *userRepository) Update(user *models.User) error {
	return ur.db.Save(user).Error
}

// Delete deletes an existing user.
func (ur *userRepository) Delete(user *models.User) error {
	return ur.db.Delete(user).Error
}

// // FindByID retrieves a user by ID.
// func (r *UserRepository) FindByID(id int) (*models.User, error) {
// 	var user models.User
// 	err := r.db.First(&user, id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// // FindByUsername retrieves a user by username.
// func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
// 	var user models.User
// 	err := r.db.Where("username = ?", username).First(&user).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// // FindByEmail retrieves a user by email.
// func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
// 	var user models.User
// 	err := r.db.Where("email = ?", email).First(&user).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// // GetByID returns a user with the given ID.
// func (ur *UserRepository) GetByID(id uint) (*models.User, error) {
// 	var user models.User
// 	result := ur.db.First(&user, id)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &user, nil
// }

// // GetByUsername returns a user with the given username.
// func (ur *UserRepository) GetByUsername(username string) (*models.User, error) {
// 	var user models.User
// 	result := ur.db.Where("username = ?", username).First(&user)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &user, nil
// }

// // GetByEmail returns a user with the given email.
// func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
// 	var user models.User
// 	result := ur.db.Where("email = ?", email).First(&user)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &user, nil
// }

// // Create creates a new user.
// func (ur *UserRepository) Create(user *models.User) error {
// 	result := ur.db.Create(user)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// // Update updates an existing user.
// func (ur *UserRepository) Update(user *models.User) error {
// 	result := ur.db.Save(user)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// // Delete deletes an existing user.
// func (ur *UserRepository) Delete(user *models.User) error {
// 	result := ur.db.Delete(user)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }
