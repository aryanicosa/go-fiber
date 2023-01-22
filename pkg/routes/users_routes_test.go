package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/repository"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"github.com/aryanicosa/go-fiber-rest-api/platform/database"
	"github.com/aryanicosa/go-fiber-rest-api/platform/migrations"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

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
		log.Fatal("fail to load database")
	}

	// migration
	migrationFileSource := os.Getenv("SQL_SOURCE_PATH")
	err = migrations.Migrate(migrationFileSource)
	if err != nil {
		log.Fatal("database migration fail")
	}

	// Define routes.
	UsersRoutes(app)
}

func TestUserSignUp(t *testing.T) {
	test := struct {
		route        string // input route
		expectedCode int
	}{
		route:        "/v1/user/sign/up",
		expectedCode: 201,
	}

	suffix := utils.String(12)
	reqBody := &models.SignUp{
		Email:    fmt.Sprintf("test%s@mail.com", suffix),
		Password: "Password123",
		UserRole: "admin",
	}
	reqBodyStr, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic YWRtaW46c2VjcmV0")

	// Perform the request plain with the app.
	resp, err := app.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign up user test")
	}

	var userSignUpResponse models.User
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &userSignUpResponse)

	defer func() {
		db, err := database.UserDB()
		if err != nil {
			fmt.Println("fail to connect user db")
		}

		err = db.DeleteUser(userSignUpResponse.ID)
		if err != nil {
			fmt.Println("fail to delete user")
		}
	}()

	fmt.Print(string(responseBodyBytes))
	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, reqBody.Email, userSignUpResponse.Email)
}

func TestUserSignIn(t *testing.T) {
	db, err := database.UserDB()
	if err != nil {
		fmt.Println("fail connect user db")
	}

	suffix := utils.String(12)
	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        fmt.Sprintf("test%s@mail.com", suffix),
		PasswordHash: utils.GeneratePassword("Password123"),
		UserStatus:   0,
		UserRole:     repository.UserRoleName,
	}
	err = db.CreateUser(user)
	if err != nil {
		log.Fatal("unable to create user")
	}

	test := struct {
		route        string // input route
		expectedCode int
	}{
		route:        "/v1/user/sign/in",
		expectedCode: 200,
	}

	reqBody := &models.SignIn{
		Email:    user.Email,
		Password: "Password123",
	}
	reqBodyStr, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic YWRtaW46c2VjcmV0")

	// Perform the request plain with the app.
	resp, err := app.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign in user test")
	}

	var userSignInResponse utils.Tokens
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &userSignInResponse)

	defer func() {
		err = db.DeleteUser(user.ID)
		if err != nil {
			fmt.Println("fail to delete user")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.NotEmpty(t, userSignInResponse.AccessToken)
	assert.NotEmpty(t, userSignInResponse.RefreshToken)
}

func TestUserRenewToken(t *testing.T) {
	db, err := database.UserDB()
	if err != nil {
		fmt.Println("fail connect user db")
	}

	suffix := utils.String(12)
	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        fmt.Sprintf("test%s@mail.com", suffix),
		PasswordHash: utils.GeneratePassword("Password123"),
		UserStatus:   0,
		UserRole:     repository.UserRoleName,
	}
	err = db.CreateUser(user)
	if err != nil {
		log.Fatal("unable to create user")
	}

	test := struct {
		route        string // input route
		expectedCode int
	}{
		route:        "/v1/user/sign/in",
		expectedCode: 200,
	}

	reqBody := &models.SignIn{
		Email:    user.Email,
		Password: "Password123",
	}
	reqBodyStr, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic YWRtaW46c2VjcmV0")

	// Perform the request plain with the app.
	resp, err := app.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign in user test")
	}

	var userSignInResponse utils.Tokens
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &userSignInResponse)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.NotEmpty(t, userSignInResponse.AccessToken)
	assert.NotEmpty(t, userSignInResponse.RefreshToken)

	testRenew := struct {
		route        string // input route
		expectedCode int
	}{
		route:        "/v1/user/token/renew",
		expectedCode: 200,
	}

	reqRenewBody := &models.Renew{
		RefreshToken: userSignInResponse.RefreshToken,
	}
	reqBodyStr, _ = json.Marshal(reqRenewBody)

	req = httptest.NewRequest("POST", testRenew.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+userSignInResponse.AccessToken)

	// Perform the request plain with the app.
	resp, err = app.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to renew user token test")
	}

	var userRenewResponse utils.Tokens
	responseBodyBytes, _ = ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &userRenewResponse)

	defer func() {
		err = db.DeleteUser(user.ID)
		if err != nil {
			fmt.Println("fail to delete user")
		}
	}()

	assert.Equal(t, testRenew.expectedCode, resp.StatusCode)
	assert.NotEmpty(t, userRenewResponse.AccessToken)
	assert.NotEmpty(t, userRenewResponse.RefreshToken)
}

func TestUserSignOut(t *testing.T) {
	tokenOnly, err := utils.GenerateNewTokens(
		uuid.New().String(),
		[]string{},
	)
	if err != nil {
		panic(err)
	}

	testSignOut := struct {
		route        string // input route
		expectedCode int
	}{
		route:        "/v1/user/sign/out",
		expectedCode: 204,
	}

	req := httptest.NewRequest("POST", testSignOut.route, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tokenOnly.AccessToken)

	// Perform the request plain with the app.
	resp, err := app.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to renew user token test")
	}

	assert.Equal(t, testSignOut.expectedCode, resp.StatusCode)
}
