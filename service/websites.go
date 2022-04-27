package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

// Website struct that represents JSON request body
type Website struct {
	// Url           string `json:"url" validate:"required,url"`
	// Name  string `json:"website" validate:"required"`
	// Expires       string `json:"expires"`
	Token         string    `json:"token" validate:"required"`
	UserId        uuid.UUID `json:"user_id" validate:"required"`
	Authenticated bool      `json:"authenticated"`
	Referer       string    `json:"referer"`
}

// Create and add a new website
func (r *Repository) CreateWebsite(context *fiber.Ctx) error {
	website := models.Websites{}

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
		// Url:           website.Url,
		// Name:  website.Name,
		Token: website.Token,
		// Expires:       website.Expires,
		UserId:        website.UserId,
		Authenticated: false,
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
			"message": "could not create website",
		})
		return err
	}

	// 200 response
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "website added successfully",
	})

	return nil
}

func (r *Repository) CreateUserWebsite(context *fiber.Ctx, token *models.AuthTokens) error {
	website := &models.Websites{}

	err := context.BodyParser(&website)
	url := fmt.Sprintf("%v", context.GetReqHeaders()["Referer"])
	// fmt.Println(context.GetReqHeaders())

	if err != nil {
		// context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
		// 	"message": "failed to process",
		// })
		return err
	}

	websiteData := &Website{
		// Url:           website.Url,
		// Name:  website.Name,
		Token:   token.Token,
		Referer: url,
		// Expires:       website.Expires,
		UserId:        token.UserId,
		Authenticated: false,
	}

	validator := validator.New()
	err = validator.Struct(websiteData)
	if err != nil {
		// context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
		// 	"message": "failed to parse JSON",
		// })
		return err
	}

	// Create row in database
	// err = r.DB.Create(&websiteData).Error
	result := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "authenticated", "referer"}),
	}).Create(&websiteData)
	err = result.Error
	if err != nil {
		// context.Status(http.StatusBadRequest).JSON(&fiber.Map{
		// 	"message": "could not create websie",
		// })
		return err
	}

	// 200 response
	// context.Status(http.StatusCreated).JSON(&fiber.Map{
	// 	"message": "website added successfully",
	// 	"id":      token.UserId,
	// })

	return nil
}

func (r *Repository) VerifyUserWebsite(context *fiber.Ctx, token *AuthToken) error {
	result := r.DB.Model(models.Websites{}).Where("token = ? AND user_id = ?", token.Token, token.UserId).Update("authenticated", "true")
	err := result.Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not update website",
		})
		return err
	}

	website := &models.Websites{}
	result = r.DB.Model(models.Websites{}).Where("token = ? AND user_id = ?", token.Token, token.UserId).First(&website)
	err = result.Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not find website",
		})
		return err
	}
	// fmt.Println(website)

	// url := fmt.Sprintf("%v", context.GetReqHeaders()["Referer"])
	authWebsite := &models.AuthedWebsites{
		Token:   website.Token,
		UserId:  website.UserId,
		Referer: website.Referer,
		// Expired: false,
		Created: datatypes.Date(time.Now()),
		Expires: datatypes.Date(time.Now().AddDate(0, 0, 7)),
	}

	result = r.DB.Model(models.AuthedWebsites{}).Create(authWebsite)
	err = result.Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not create authed website",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user has been successfully authenticated",
		"id":      token.UserId,
		// "email":   .Email,
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
		// Url:     website.Url,
		// Expires: website.Expires,
		// Name:   website.Name,
		Token:  website.Token,
		UserId: website.UserId,
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
