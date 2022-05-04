package service

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
)

func (r *Repository) Login(context *fiber.Ctx) error {
	users := &models.Users{}

	// String value of the id section of the url
	username := context.Params("username")
	if len(username) < 1 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "username cannot have less than 1 characters",
			// "data":    users,
		})
	} else {
		// Select from users where the id is the same as the one passed through the url
		result := r.DB.Where("username = ?", username).Find(&users)
		err := result.Error
		if err != nil {
			context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not fetch user"})
			return err
		}

		if result.RowsAffected > 0 {
			// fmt.Println("Found")
			r.CreateMagicToken(users.ID, users.Email, context)

			context.Status(http.StatusOK).JSON(&fiber.Map{
				"message":  "user retrieved successfully",
				"id":       users.ID,
				"username": users.Username,
				"email":    users.Email,
			})
		} else {
			// Return the users in the JSON response
			context.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "user could not be retrieved successfully",
				// "data":    users,
			})
		}
	}

	return nil
}

func (r *Repository) generateQRCode() (*os.File, string) {
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
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
	// defer os.Remove(file.Name())
	png.Encode(file, qrCode)

	return file, stringToSHA512
}

// Generate the QR code for a random 64 byte set of characters
func (r *Repository) CreateQRCode(context *fiber.Ctx) error {
	id := context.Params("id")
	url := fmt.Sprintf("%v", context.GetReqHeaders()["Referer"])

	// fmt.Println(url)
	qrImage, qrCode := r.generateQRCode()
	defer os.Remove(qrImage.Name())

	qrData := &AuthToken{
		Token:  qrCode,
		UserId: uuid.Must(uuid.FromString(id)),
		// Authenticated: false,
	}

	result := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token"}),
	}).Create(&qrData)
	err := result.Error
	if err != nil {
		// log.Fatal(err)
		return err
	}

	websiteData := &Website{
		// Name:          website.Name,
		Referer:       url,
		Token:         qrCode,
		UserId:        uuid.Must(uuid.FromString(id)),
		Authenticated: false,
	}

	validator := validator.New()
	err = validator.Struct(websiteData)
	if err != nil {
		return err
	}

	// Create row in database
	result = r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "authenticated", "referer"}),
	}).Create(&websiteData)

	err = result.Error
	if err != nil {
		return err
	}

	// Return the file as a response
	context.Status(http.StatusOK).SendFile(qrImage.Name(), true)

	// r.CheckUserVerified(context)

	return nil
}

func (r *Repository) VerifyQRCode(context *fiber.Ctx) error {
	qrCode := &models.AuthTokens{}

	err := context.BodyParser(&qrCode)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Could not process body",
		})
	}

	qrCodeData := &AuthToken{
		Token:  qrCode.Token,
		UserId: qrCode.UserId,
	}

	validator := validator.New()
	validator.Struct(qrCodeData)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to parse JSON",
		})
	}

	result := r.DB.Model(models.AuthTokens{}).Where("user_id = ? AND token = ?", qrCode.UserId, qrCode.Token).Find(&qrCode)
	err = result.Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "invalid token"})
		return err
	}

	if result.RowsAffected > 0 {
		r.VerifyUserWebsite(context, qrCodeData)

		context.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "success",
			"id":      qrCode.ID,
			"user_id": qrCode.UserId,
		})
	} else {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "invalid qrCode",
		})
	}

	return nil
}

func (r *Repository) CheckUserVerified(context *fiber.Ctx) error {
	// fiber.AcquireAgent().Timeout(time.Second * 30)
	qrCode := &models.AuthTokens{}
	url := fmt.Sprintf("%v", context.GetReqHeaders()["Referer"])

	id := context.Params("id")

	success := false

	for !success {
		// fmt.Println(success)
		result := r.DB.Model(models.Websites{}).Where("user_id = ? AND referer = ? AND authenticated = 'true'", id, url).Find(&qrCode)
		err := result.Error
		if err != nil {
			context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "invalid token"})
			return err
		}

		if result.RowsAffected > 0 {
			success = true
			user := &models.Users{}

			result = r.DB.Model(user).Where("id = ?", id).Find(&user)
			err = result.Error
			if err != nil {
				context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "user not found"})
				return err
			}

			context.Status(http.StatusOK).JSON(&fiber.Map{
				// "message": "success",
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
			})
		} else {
			success = false
		}
	}

	// context.Status(http.StatusGatewayTimeout).JSON(&fiber.Map{
	// 	"message": "failed to verify user",
	// })

	return nil
}

func (r *Repository) GetLatestToken(context *fiber.Ctx) error {
	token := &models.AuthTokens{}
	id := context.Params("id")

	result := r.DB.Model(models.AuthTokens{}).Where("user_id = ?", id).Find(&token)
	err := result.Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not fetch user token"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"token": token.Token,
		// "referer": website.Referer,
		"userId": token.UserId,
	})

	return nil
}
