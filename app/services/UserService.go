package services

import (
	"errors"

	"github.com/mousav1/weiser/app/models"
	"github.com/mousav1/weiser/app/repositories"
)

// UserService represents the service for managing users.
type UserService interface {
	CreateUser(string, string, string) (*models.User, error)
	GetUserByID(uint) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	UpdateUser(uint, string, string, string) error
	DeleteUser(uint) error
}

// UpdateUserInput represents the input required for updating a user.
type UpdateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=6,max=20"`
}

// CreateUserInput represents the input required for creating a user.
type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userService struct {
	userRepository repositories.UserRepository
}

// NewUserService creates a new instance of userService.
func NewUserService(ur repositories.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

// CreateUser creates a new user.
func (us *userService) CreateUser(username, email, password string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	err := us.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by its ID.
func (us *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := us.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByUsername retrieves a user by its username.
func (us *userService) GetUserByUsername(username string) (*models.User, error) {
	user, err := us.userRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByEmail retrieves a user by its email.
func (us *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := us.userRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// UpdateUser updates an existing user.
func (us *userService) UpdateUser(id uint, username, email, password string) error {
	user, err := us.GetUserByID(id)
	if err != nil {
		return err
	}
	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		user.Password = password
	}
	return us.userRepository.Update(user)
}

// DeleteUser deletes an existing user.
func (us *userService) DeleteUser(id uint) error {
	user, err := us.GetUserByID(id)
	if err != nil {
		return err
	}
	return us.userRepository.Delete(user)
}

// // GetByID returns a user with the given ID.
// func (us *userService) GetByID(id uint) (*models.User, error) {
// 	return us.userRepository.GetByID(id)
// }

// // GetByUsername returns a user with the given username.
// func (us *userService) GetByUsername(username string) (*models.User, error) {
// 	return us.userRepository.GetByUsername(username)
// }

// // GetByEmail returns a user with the given email.
// func (us *userService) GetByEmail(email string) (*models.User, error) {
// 	return us.userRepository.GetByEmail(email)
// }

// // Create creates a new user.
// func (us *userService) Create(user *models.User) error {
// 	return us.userRepository.Create(user)
// }

// // Update updates an existing user.
// func (us *userService) Update(user *models.User) error {
// 	return us.userRepository.Update(user)
// }

// // Delete deletes an existing user.
// func (us *userService) Delete(user *models.User) error {
// 	return us.userRepository.Delete(user)
// }
