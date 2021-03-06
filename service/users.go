package service

import (
	"net/http"
	"time"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

// The user data sent to the server
type User struct {
	Username string         `json:"username" validate:"required"`
	Email    string         `json:"email" validate:"required"`
	Created  datatypes.Date `json:"created"`
	// Website  string `json:"website" validate:"required"`
}

// Create and add a new user to the database
func (r *Repository) CreateUser(context *fiber.Ctx) error {
	user := models.Users{
		Created: datatypes.Date(time.Now()),
	}

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
		Created:  user.Created,
		// Website:  user.Website,
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

	result := r.DB.Model(models.Users{}).Create(&user)
	err = result.Error
	// Create a new user in database
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not create user",
		})
		return err
	}

	// Return 200 response
	context.Status(http.StatusCreated).JSON(&fiber.Map{
		"message":  "user has been successfully created",
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})

	r.CreateMagicToken(user.ID, user.Email, context)

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
	user := &models.Users{}

	// String value of the id section of the url
	id := context.Params("id")

	// Select from users where the id is the same as the one passed through the url
	err := r.DB.Where("id =?", id).Find(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not fetch user"})
		return err
	}

	// Return the users in the JSON response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user retrieved successfully",
		// "data":    user,
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"created":       user.Created,
		"tokenDuration": user.TokenDuration,
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
	// userData := &User{
	// 	Username: user.Username,
	// 	Email:    user.Email,
	// }

	// Validate the JSON to verify integrity
	// validator := validator.New()
	// err = validator.Struct(userData)
	// if err != nil {
	// 	context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
	// 		"message": "failed to parse JSON",
	// 	})
	// 	return err
	// }

	// Update specific user using the Users model and the id passed in
	err = r.DB.Model(models.Users{}).Where("id = ?", id).Updates(&user).Error
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
