package service

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/dchest/uniuri"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	mail "github.com/xhit/go-simple-mail/v2"
)

type AuthToken struct {
	Token  string    `json:"token" validate:"required"`
	UserId uuid.UUID `json:"user_id" validate:"required"`
}

func (r *Repository) CreateMagicToken(userId uuid.UUID, email string, context *fiber.Ctx) error {
	token := uniuri.NewLen(6)
	tokenModel := &models.AuthTokens{UserId: userId, Token: token}

	result := r.DB.Create(&tokenModel)
	err := result.Error
	if err != nil {
		return err
	}

	sendEmail(token, email)

	r.CreateUserWebsite(context, tokenModel)

	return nil
}

func (r *Repository) VerifyMagicToken(context *fiber.Ctx) error {
	token := &models.AuthTokens{}

	err := context.BodyParser(&token)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Could not process body",
		})
	}

	tokenData := &AuthToken{
		Token:  token.Token,
		UserId: token.UserId,
	}

	validator := validator.New()
	err = validator.Struct(tokenData)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "failed to parse JSON",
		})
	}

	result := r.DB.Model(models.AuthTokens{}).Where("user_id = ? AND token = ?", token.UserId, token.Token).Find(&token)
	err = result.Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "invalid token"})
		return err
	}

	if result.RowsAffected > 0 {
		r.VerifyUserWebsite(context, tokenData)

		context.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "success",
			"data":    token.ID,
		})
	} else {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "invalid token",
		})
	}

	return nil
}

func sendEmail(token, userEmail string) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env missing in environment")
		return err
	}

	// from := os.Getenv("EMAIL_USERNAME")
	// password := os.Getenv("EMAIL_PASSWORD")
	// to := []string{email}
	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"

	// message := []byte("My super secret message is this token\n" + token)

	// // Create authentication
	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// // Send actual message
	// err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	// if err != nil {
	// 	log.Panic(err)
	// 	return err
	// }

	var htmlBody = fmt.Sprintf(`
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<title>Hello, World</title>
</head>

<body>
	<h3>Welcome to Direct Security</h3>
	<p>Your Security Token is: <span>%v</span></p>
</body>
</html>
	`, token)

	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = os.Getenv("EMAIL_USERNAME")
	server.Password = os.Getenv("EMAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("Direct Security <salimabdu008@gmail.com>")
	email.AddTo(userEmail)
	email.SetSubject("Direct Security Authentication Token")
	email.SetBody(mail.TextHTML, htmlBody)

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}
