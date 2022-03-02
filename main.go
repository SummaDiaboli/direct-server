package main

import (
	"log"
	"os"

	"github.com/SummaDiaboli/direct-server/models"
	"github.com/SummaDiaboli/direct-server/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Retreive environment variables
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Instantiate database struct
	// config := &database.Config{
	// 	Host:     os.Getenv("DB_HOST"),
	// 	Port:     os.Getenv("DB_PORT"),
	// 	Password: os.Getenv("DB_PASS"),
	// 	User:     os.Getenv("DB_USER"),
	// 	SSLMode:  os.Getenv("DB_SSLMODE"),
	// 	DBName:   os.Getenv("DB_NAME"),
	// }

	dbURL := os.Getenv("DATABASE_URL")

	// Create a new database connection
	// db, err := database.NewConnection(config)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("could not load database")
	}

	// Migrate database tables
	err = models.MigrateTables(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	// Database instance
	r := &service.Repository{
		DB: db,
	}

	// Start fiber server
	app := fiber.New()
	app.Use(logger.New())

	r.SetupRoutes(app)
	app.Listen(":8080")
}
