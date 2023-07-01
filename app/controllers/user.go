package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/http/request"
	"github.com/mousav1/weiser/app/http/response"

	"github.com/mousav1/weiser/app/services"
)

// UserController represents the controller for managing users.
type UserController interface {
	CreateUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	GetUserByUsername(c *fiber.Ctx) error
	GetUserByEmail(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

// NewUserController creates a new instance of userController.
func NewUserController(us services.UserService) UserController {
	return &userController{
		userService: us,
	}
}

// CreateUser creates a new user.
func (uc *userController) CreateUser(c *fiber.Ctx) error {
	req, err := request.New(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user services.CreateUserInput
	if err := req.Bind(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createdUser, err := uc.userService.CreateUser(user.Username, user.Email, user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": createdUser.ID,
	})
}

// GetUserByID retrieves a user by its ID.
func (uc *userController) GetUserByID(c *fiber.Ctx) error {
	req, err := request.New(c)
	if err != nil {
		res := response.New(nil, err.Error(), fiber.StatusInternalServerError)
		return response.Send(c, res)
	}

	id, err := req.Int("id")
	if err != nil {
		res := response.New(nil, "Invalid user ID", fiber.StatusBadRequest)
		return response.Send(c, res)
	}

	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			res := response.New(nil, ErrUserNotFound.Error(), fiber.StatusNotFound)
			return response.Send(c, res)
		}
		res := response.New(nil, err.Error(), fiber.StatusInternalServerError)
		return response.Send(c, res)
	}

	res := response.New(user, "User data retrieved successfully", fiber.StatusOK)
	return response.Send(c, res)
}

// GetUserByUsername retrieves a user by its username.
func (uc *userController) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := uc.userService.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": ErrUserNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// GetUserByEmail retrieves a user by its email.
func (uc *userController) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	user, err := uc.userService.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": ErrUserNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateUser updates an existing user.
func (uc *userController) UpdateUser(c *fiber.Ctx) error {
	req, err := request.New(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user services.UpdateUserInput
	if err := req.Bind(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	idString := c.Params("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = uc.userService.UpdateUser(uint(id), user.Username, user.Email, user.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": ErrUserNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// data := views.ViewData{
	// 	Title: "Home",
	// 	Data:  "John Smith",
	// }
	// err := views.View(c, data, "test.html")
	// if err != nil {
	// 	return err
	// }
	// return nil

	return c.SendStatus(fiber.StatusNoContent)

}

// DeleteUser deletes an existing user.
func (uc *userController) DeleteUser(c *fiber.Ctx) error {
	idString := c.Params("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = uc.userService.DeleteUser(uint(id))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": ErrUserNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)

}

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")
