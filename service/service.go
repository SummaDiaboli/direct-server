package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
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
	api.Get("/resend-token/:id", r.ResendMagicToken)

	// Authentication API
	api.Get("/login/:username", r.Login)
	api.Post("/confirm-qr", r.VerifyQRCode)
	api.Get("/verify/:id", r.CheckUserVerified)

	// Authed Websites
	api.Get("/authed-websites/:id", r.GetAuthedWebsites)

	// Util API
	api.Get("/generate/:id", timeout.New(r.CreateQRCode, 30*time.Second))
	// api.Get("/", sayHi)
}
