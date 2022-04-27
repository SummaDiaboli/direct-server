package service

import (
	"net/http"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateAuthedWebsite(context *fiber.Ctx) error {
	return nil
}

func (r *Repository) GetAuthedWebsites(context *fiber.Ctx) error {
	// authedWebsite := &models.AuthedWebsites{}

	id := context.Params("id")

	results := []models.AuthedWebsites{}

	result := r.DB.Model(&models.AuthedWebsites{}).Where("user_id = ?", id).Find(&results)
	err := result.Error
	// err := r.DB.Where("id =?", id).Find(&users).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not fetch user"})
		return err
	}

	// fmt.Println(result.RowsAffected)

	// Return the users in the JSON response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":  "user retrieved successfully",
		"websites": results,
	})

	return nil
}

func (r *Repository) GetAuthedWebsiteById(context *fiber.Ctx) error {
	return nil
}

func (r *Repository) UpdateAuthedWebsite(context *fiber.Ctx) error {
	return nil
}

func (r *Repository) DeleteAuthedWebsite(context *fiber.Ctx) error {
	return nil
}
