package service

import (
	"net/http"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Website struct that represents JSON request body
type Website struct {
	Url       string `json:"url" validate:"required,url"`
	Name      string `json:"name" validate:"required"`
	UserToken string `json:"user_token" validate:"required"`
	Expires   string `json:"expires"`
}

// Create and add a new website
func (r *Repository) CreateWebsite(context *fiber.Ctx) error {
	website := Website{}

	// Parse body JSON into website struct
	err := context.BodyParser(&website)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to process",
		})
		return err
	}

	// Struct for validation
	websiteData := &Website{
		Url:       website.Url,
		Name:      website.Name,
		UserToken: website.UserToken,
		Expires:   website.Expires,
	}

	// Verify that JSON integrity
	validator := validator.New()
	err = validator.Struct(websiteData)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to parse JSON",
		})
		return err
	}

	// Create row in database
	err = r.DB.Create(&website).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not create websie",
		})
		return err
	}

	// 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "website added successfully",
	})

	return nil
}

// Get all website in table
func (r *Repository) GetWebsites(context *fiber.Ctx) error {
	websites := &[]models.Websites{}

	// Select and return all websites in table
	err := r.DB.Find(websites).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not fetch websites"},
		)
		return err
	}

	// 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "websites have been successfully retrieved",
		"data":    websites,
	})

	return nil
}

// Update specific website in table
func (r *Repository) UpdateWebsite(context *fiber.Ctx) error {
	website := &models.Websites{}

	// id parameter in url path
	id := context.Params("id")

	// Parse body JSON into website struct
	err := context.BodyParser(&website)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to process",
		})
		return err
	}

	// Map body to struct for validation
	websiteData := &Website{
		Url:       website.Url,
		Expires:   website.Expires,
		Name:      website.Name,
		UserToken: website.UserToken,
	}

	// Verify website struct integrity
	validator := validator.New()
	err = validator.Struct(websiteData)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to parse JSON",
		})
		return err
	}

	// Update specific row in table where id matches using the website struct
	err = r.DB.Model(models.Websites{}).Where("id = ?", id).Updates(websiteData).Error
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "could not update website",
		})
		return err
	}

	// 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "website has been successfully updated",
	})

	return nil
}

// Delete specific website
func (r *Repository) DeleteWebsite(context *fiber.Ctx) error {
	website := &models.Websites{}

	// Id parameter in url path
	id := context.Params("id")

	// Delete row in table where website id matches parameter id
	err := r.DB.Where("id = ?", id).Delete(&website).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not fetch website"})
		return err
	}

	// 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "website deleted successfully",
	})

	return nil
}
