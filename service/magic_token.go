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
	"gorm.io/gorm/clause"
)

type AuthToken struct {
	Token  string    `json:"token" validate:"required"`
	UserId uuid.UUID `json:"user_id" validate:"required"`
}

func (r *Repository) CreateMagicToken(userId uuid.UUID, email string, context *fiber.Ctx) error {
	token := uniuri.NewLen(6)
	tokenModel := &models.AuthTokens{UserId: userId, Token: token}

	result := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token"}),
	}).Create(&tokenModel)
	err := result.Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	context.Append("user_id", userId.String())

	// fmt.Println("created token")
	// log.Println("token created")
	sendEmail(token, email)

	r.CreateUserWebsite(context, tokenModel)

	return nil
}

func (r *Repository) VerifyMagicToken(context *fiber.Ctx) error {
	token := &models.AuthTokens{}
	user := &models.Users{}

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

		result = r.DB.Model(user).Where("id = ?", token.UserId).Find(&user)
		err = result.Error
		if err != nil {
			return err
		}

		if result.RowsAffected > 0 {
			context.Status(http.StatusOK).JSON(&fiber.Map{
				"message":  "user retrieved successfully",
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			})
		} else {
			// Return the users in the JSON response
			context.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"message": "user could not be retrieved successfully",
				// "data":    users,
			})
		}
	} else {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "invalid token",
		})
	}

	return nil
}

func (r *Repository) ResendMagicToken(context *fiber.Ctx) error {
	token := uniuri.NewLen(6)
	user := models.Users{}
	userId := context.Params("id")

	tokenModel := &models.AuthTokens{UserId: uuid.FromStringOrNil(userId), Token: token}
	result := r.DB.Model(models.AuthTokens{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token"}),
	}).Create(&tokenModel)
	err := result.Error
	if err != nil {
		return err
	}

	result = r.DB.Model(models.Users{}).Where("id = ?", userId).Find(&user)
	err = result.Error
	if err != nil {
		return err
	}

	sendEmail(token, user.Email)

	r.CreateUserWebsite(context, tokenModel)

	context.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": "magic code recreated",
		"id":      userId,
	})

	return nil
}

// TODO: Implement Resend email service
// TODO: Resend email should fetch token from referer and user_id, then send it to email
/* TODO: Resend email */
func sendEmail(token, userEmail string) error {
	// log.Println("Here")
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
		log.Println(err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("Direct Security <salimabdu008@gmail.com>")
	email.AddTo(userEmail)
	email.SetSubject("Direct Security Authentication Token")
	email.SetBody(mail.TextHTML, htmlBody)

	err = email.Send(smtpClient)
	// log.Println("email sent")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
