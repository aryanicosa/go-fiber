package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"github.com/aryanicosa/go-fiber-rest-api/platform/database"
	"github.com/aryanicosa/go-fiber-rest-api/platform/migrations"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

func TestUserRoutes(t *testing.T) {
	// Load .env.test file from the root folder
	if err := godotenv.Load("../../.env.test"); err != nil {
		panic(err)
	}

	// Define Fiber app.
	app = fiber.New()

	// init connect to db
	_, err := database.InitDBConnection()
	if err != nil {
		log.Fatal("could not load database")
	}

	// migration
	migrationFileSource := os.Getenv("SQL_SOURCE_PATH")
	err = migrations.Migrate(migrationFileSource)
	if err != nil {
		log.Fatal("database migration fail")
	}

	// Define routes.
	UsersRoutes(app)

	TestUserSignUp(t)
}

func TestUserSignUp(t *testing.T) {
	test := struct {
		description  string
		route        string // input route
		expectedCode int
	}{
		description:  "user sign up",
		route:        "/v1/user/sign/up",
		expectedCode: 201,
	}

	suffix := utils.String(12)
	reqBody := &models.SignUp{
		Email:    fmt.Sprintf("test%s@mail.com", suffix),
		Password: "Password123",
		UserRole: "admin",
	}
	resBodyStr, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", test.route, bytes.NewBufferString(string(resBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	// Perform the request plain with the app.
	resp, err := app.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail sign up user test")
	}

	var userSignUpResponse models.User
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &userSignUpResponse)

	fmt.Println(resp)
	fmt.Println(resp.Body)
	fmt.Println(string(responseBodyBytes))
	fmt.Println(userSignUpResponse)

	db, err := database.UserDB()
	if err != nil {
		fmt.Println("fail connect user db")
	}
	err = db.DeleteUser(userSignUpResponse.ID)
	if err != nil {
		fmt.Println("fail delete user")
	}

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, reqBody.Email, userSignUpResponse.Email)
}

func TestUserSignIn(t *testing.T) {

}

func TestUserRenewToken(t *testing.T) {

}

func TestUserSignOut(t *testing.T) {

}
