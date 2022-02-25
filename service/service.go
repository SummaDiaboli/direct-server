package service

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image/png"
	"net/http"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// The database in use
type Repository struct {
	DB *gorm.DB
}

// Generate the QR code for a random 64 byte set of characters
func generateQR(context *fiber.Ctx) error {

	// Create a 64 byte key using the rand package
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "key generation went wrong",
			},
		)
	}

	// Hash the key using SHA512
	hasher := crypto.SHA512.New()
	hasher.Write(key)
	stringToSHA512 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Generate the QR code
	qrCode, _ := qr.Encode(fmt.Sprintf("%v", stringToSHA512), qr.H, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 300, 300)

	// Save the QR code as an image, with random characters appended to it to prevent overriting
	file, _ := os.CreateTemp("", "qrcode*.png")
	defer file.Close()
	defer os.Remove(file.Name())
	png.Encode(file, qrCode)

	// Return the file as a response
	context.Status(http.StatusOK).SendFile(file.Name())

	return nil
}

func sayHi(context *fiber.Ctx) error {
	context.Status(http.StatusOK).SendString("Hello, world!")

	return nil
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

	// Util API
	api.Get("/generate", generateQR)
	api.Get("/", sayHi)
}
