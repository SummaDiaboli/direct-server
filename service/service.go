package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// The database in use
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// User API
	api.Post("/users", r.CreateUser)
	api.Get("/users", r.GetUsers)
	api.Get("/users/:id", r.GetUserById)
	api.Delete("/users/:id", r.DeleteUser)
	api.Patch("/users/:id", r.UpdateUser)

	// Website API
	api.Post("/websites", r.CreateWebsite)
	api.Get("/websites", r.GetWebsites)
	api.Delete("/websites/:id", r.DeleteWebsite)
	api.Patch("/websites/:id", r.UpdateWebsite)

	// Magic Token
	api.Post("/confirm-token", r.VerifyMagicToken)

	// Authentication API
	api.Get("/login/:username", r.Login)
	api.Post("/confirm-qr", r.VerifyQRCode)

	// Util API
	api.Get("/generate/:id", r.CreateQRCode)
	// api.Get("/", sayHi)
}
