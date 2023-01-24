package routes

import (
	"github.com/gofiber/fiber/v2"
)

var AppTest *fiber.App

//func TestMain(m *testing.M) {
//	// Load .env.test file from the root folder.
//	if err := godotenv.Load("../../.env.test"); err != nil {
//		log.Fatal(err)
//	}
//
//	// Define Fiber AppTest.
//	AppTest = fiber.New()
//
//	// init connect to db
//	_, err := database.InitDBConnection()
//	if err != nil {
//		log.Fatal("fail to load database")
//	}
//
//	// migration
//	migrationFileSource := os.Getenv("SQL_SOURCE_PATH")
//	err = migrations.Migrate(migrationFileSource)
//	if err != nil {
//		log.Fatal("database migration fail")
//	}
//
//	// Define routes.
//	BooksRoutes(AppTest)
//	UsersRoutes(AppTest)
//	MiscRoutes(AppTest)
//
//	os.Exit(m.Run())
//}
