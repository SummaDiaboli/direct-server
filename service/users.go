package service

import (
	"net/http"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// The user data sent to the server
type User struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

// Create and add a new user to the database
func (r *Repository) CreateUser(context *fiber.Ctx) error {
	user := User{}

	// Parse the JSON body of the request
	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to process",
		})
		return err
	}

	// Model for data validation
	userData := &User{
		Username: user.Username,
		Email:    user.Email,
	}

	// Validate the JSON to verify integrity
	validator := validator.New()
	err = validator.Struct(userData)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to parse JSON",
		})
		return err
	}

	// Create a new user in database
	err = r.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not create user",
		})
		return err
	}

	// Return 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user has been successfully created",
	})

	return nil
}

// Get all the users in the database
func (r *Repository) GetUsers(context *fiber.Ctx) error {
	users := &[]models.Users{}

	// Select and return all users in the database
	err := r.DB.Find(users).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not fetch users"},
		)
		return err
	}

	// Return 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "users have been successfully retrieved",
		"data":    users,
	})

	return nil
}

// Return user with specific ID
func (r *Repository) GetUserById(context *fiber.Ctx) error {
	users := &models.Users{}

	// String value of the id section of the url
	id := context.Params("id")

	// Select from users where the id is the same as the one passed through the url
	err := r.DB.Where("id =?", id).Find(&users).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not fetch user"})
		return err
	}

	// Return the users in the JSON response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user retrieved successfully",
		"data":    users,
	})

	return nil
}

// Delete specific user from database
func (r *Repository) DeleteUser(context *fiber.Ctx) error {
	users := &models.Users{}

	// String value from the url
	id := context.Params("id")

	// Deletes the user where the ID equals the value passed in
	err := r.DB.Where("id = ?", id).Delete(&users).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not fetch user"})
		return err
	}

	// 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user deleted successfully",
	})

	return nil
}

// Update specific user
func (r *Repository) UpdateUser(context *fiber.Ctx) error {
	user := &models.Users{}

	// String value from the url
	id := context.Params("id")

	// Updates the user where the ID equals the value passed in
	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to process",
		})
		return err
	}

	// User model for validation
	userData := &User{
		Username: user.Username,
		Email:    user.Email,
	}

	// Validate the JSON to verify integrity
	validator := validator.New()
	err = validator.Struct(userData)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to parse JSON",
		})
		return err
	}

	// Update specific user using the Users model and the id passed in
	err = r.DB.Model(models.Users{}).Where("id = ?", id).Updates(userData).Error
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "could not update user",
		})
		return err
	}

	// 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user has been successfully updated",
	})

	return nil
}
