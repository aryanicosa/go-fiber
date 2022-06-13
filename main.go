package main

import (
	"log"
	"os"

	"github.com/aryanicosa/go-fiber-rest-api/database"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/user"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/user/model"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatal("could not load database")
	}

	err = model.MigrateUsers(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := &user.Service{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)

	app.Listen(":8080")
}
