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
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestBookInit(t *testing.T) {
	// Load .env.test file from the root folder.
	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Fatal(err)
	}

	// Define Fiber AppTest.
	AppTest = fiber.New()

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
	BooksRoutes(AppTest)
}

func TestCreateBook(t *testing.T) {
	db := database.UserDB()

	suffix := utils.String(12)
	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        fmt.Sprintf("test%s@mail.com", suffix),
		PasswordHash: utils.GeneratePassword("Password123"),
		UserStatus:   0,
		UserRole:     repository.AdminRoleName,
	}
	err := db.CreateUser(user)
	if err != nil {
		log.Fatal("unable to create user")
	}

	// Create token with `book:create` credential.
	tokenOnlyCreate, err := utils.GenerateNewTokens(
		user.ID.String(),
		[]string{"book:create"},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define a structure for specifying input and output data of a single test case.
	test := struct {
		route        string // input route
		method       string // input method
		expectedCode int
	}{
		route:        "/v1/book",
		method:       "POST",
		expectedCode: 201,
	}

	book := &models.Book{
		Title:      "Test Title",
		Author:     "John Doe",
		BookStatus: 0,
		BookAttrs: models.BookAttrs{
			Picture:     "Test Pic",
			Description: "This book is test book",
			Rating:      6,
		},
	}
	reqBodyStr, _ := json.Marshal(book)

	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Authorization", "Bearer "+tokenOnlyCreate.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	// Perform the request plain with the AppTest.
	resp, err := AppTest.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign in user test")
	}

	var createBookResponse models.Book
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &createBookResponse)

	defer func() {
		err = db.DeleteUser(user.ID)
		if err != nil {
			log.Fatal("fail to delete user")
		}
		dbBook := database.BookDB()
		if err != nil {
			log.Fatal("fail connect book db")
		}
		err = dbBook.DeleteBook(createBookResponse.ID)
		if err != nil {
			log.Fatal("Fail to delete book")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.NotEmpty(t, createBookResponse.ID)
}

func TestGetBookById(t *testing.T) {
	db := database.UserDB()

	suffix := utils.String(12)
	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        fmt.Sprintf("test%s@mail.com", suffix),
		PasswordHash: utils.GeneratePassword("Password123"),
		UserStatus:   0,
		UserRole:     repository.AdminRoleName,
	}
	err := db.CreateUser(user)
	if err != nil {
		log.Fatal("unable to create user")
	}

	dbBook := database.BookDB()
	if err != nil {
		log.Fatal("fail connect book db")
	}

	book := &models.Book{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserID:     user.ID,
		Title:      "Test Title",
		Author:     "John Doe",
		BookStatus: 0,
		BookAttrs: models.BookAttrs{
			Picture:     "Picture",
			Description: "Description",
			Rating:      10,
		},
	}

	err = dbBook.CreateBook(book)
	if err != nil {
		log.Fatal("fail to create book")
	}

	// Define a structure for specifying input and output data of a single test case.
	test := struct {
		route        string // input route
		method       string // input method
		expectedCode int
	}{
		route:        "/v1/book/" + book.ID.String(),
		method:       "GET",
		expectedCode: 200,
	}

	req := httptest.NewRequest(test.method, test.route, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic YWRtaW46c2VjcmV0")

	// Perform the request plain with the AppTest.
	resp, err := AppTest.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign in user test")
	}

	var getBookResponse models.Book
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &getBookResponse)

	defer func() {
		err = db.DeleteUser(user.ID)
		if err != nil {
			log.Fatal("fail to delete user")
		}
		err = dbBook.DeleteBook(getBookResponse.ID)
		if err != nil {
			log.Fatal("Fail to delete book")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.NotEmpty(t, getBookResponse.ID)
}

func TestGetBookAll(t *testing.T) {
	db := database.UserDB()

	suffix := utils.String(12)
	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        fmt.Sprintf("test%s@mail.com", suffix),
		PasswordHash: utils.GeneratePassword("Password123"),
		UserStatus:   0,
		UserRole:     repository.AdminRoleName,
	}
	err := db.CreateUser(user)
	if err != nil {
		log.Fatal("unable to create user")
	}

	dbBook := database.BookDB()
	if err != nil {
		log.Fatal("fail connect book db")
	}

	books := []models.Book{
		{
			ID:         uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			UserID:     user.ID,
			Title:      "Test Title",
			Author:     "John Doe",
			BookStatus: 0,
			BookAttrs: models.BookAttrs{
				Picture:     "Picture",
				Description: "Description",
				Rating:      10,
			},
		},
		{
			ID:         uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			UserID:     user.ID,
			Title:      "Test Title 2",
			Author:     "John Doe",
			BookStatus: 0,
			BookAttrs: models.BookAttrs{
				Picture:     "Picture 2",
				Description: "Description 2",
				Rating:      10,
			},
		},
	}

	for _, book := range books {
		err = dbBook.CreateBook(&book)
		if err != nil {
			log.Fatal("fail to create book")
		}
	}

	// Define a structure for specifying input and output data of a single test case.
	test := struct {
		route        string // input route
		method       string // input method
		expectedCode int
	}{
		route:        "/v1/books",
		method:       "GET",
		expectedCode: 200,
	}

	req := httptest.NewRequest(test.method, test.route, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic YWRtaW46c2VjcmV0")

	// Perform the request plain with the AppTest.
	resp, err := AppTest.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign in user test")
	}

	var getBooksResponse []models.BookForPublic
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &getBooksResponse)

	defer func() {
		for _, book := range books {
			err = dbBook.DeleteBook(book.ID)
			if err != nil {
				log.Fatal("Fail to delete book")
			}
		}
		err = db.DeleteUser(user.ID)
		if err != nil {
			log.Fatal("fail to delete user")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	for _, bookResponse := range getBooksResponse {
		assert.NotEmpty(t, bookResponse.ID)
	}
}

func TestUpdateBookById(t *testing.T) {
	db := database.UserDB()

	suffix := utils.String(12)
	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        fmt.Sprintf("test%s@mail.com", suffix),
		PasswordHash: utils.GeneratePassword("Password123"),
		UserStatus:   0,
		UserRole:     repository.AdminRoleName,
	}
	err := db.CreateUser(user)
	if err != nil {
		log.Fatal("unable to create user")
	}

	dbBook := database.BookDB()
	if err != nil {
		log.Fatal("fail connect book db")
	}

	book := &models.Book{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserID:     user.ID,
		Title:      "Test Title",
		Author:     "John Doe",
		BookStatus: 0,
		BookAttrs: models.BookAttrs{
			Picture:     "Picture",
			Description: "Description",
			Rating:      10,
		},
	}

	err = dbBook.CreateBook(book)
	if err != nil {
		log.Fatal("fail to create book")
	}

	// Create token with `book:create` credential.
	tokenAdmin, err := utils.GenerateNewTokens(
		user.ID.String(),
		[]string{"book:create", "book:update", "book:delete"},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define a structure for specifying input and output data of a single test case.
	test := struct {
		route        string // input route
		method       string // input method
		expectedCode int
	}{
		route:        "/v1/book/" + book.ID.String(),
		method:       "PUT",
		expectedCode: 201,
	}

	bookUpdate := &models.Book{
		ID:         book.ID,
		UserID:     user.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Title:      "Test Title Update",
		Author:     "John Doe",
		BookStatus: 1,
		BookAttrs: models.BookAttrs{
			Picture:     "Test Pic update",
			Description: "This book is test book update",
			Rating:      6,
		},
	}
	reqBodyStr, _ := json.Marshal(bookUpdate)

	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Authorization", "Bearer "+tokenAdmin.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	// Perform the request plain with the AppTest.
	resp, err := AppTest.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign in user test")
	}

	var updateBookResponse models.Book
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &updateBookResponse)

	defer func() {
		err = db.DeleteUser(user.ID)
		if err != nil {
			log.Fatal("fail to delete user")
		}
		err = dbBook.DeleteBook(updateBookResponse.ID)
		if err != nil {
			log.Fatal("Fail to delete book")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.NotEmpty(t, updateBookResponse.ID)
}

func TestDeleteBookById(t *testing.T) {
	db := database.UserDB()

	suffix := utils.String(12)
	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        fmt.Sprintf("test%s@mail.com", suffix),
		PasswordHash: utils.GeneratePassword("Password123"),
		UserStatus:   0,
		UserRole:     repository.AdminRoleName,
	}
	err := db.CreateUser(user)
	if err != nil {
		log.Fatal("unable to create user")
	}

	dbBook := database.BookDB()
	if err != nil {
		log.Fatal("fail connect book db")
	}

	book := &models.Book{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserID:     user.ID,
		Title:      "Test Title",
		Author:     "John Doe",
		BookStatus: 0,
		BookAttrs: models.BookAttrs{
			Picture:     "Picture",
			Description: "Description",
			Rating:      10,
		},
	}

	err = dbBook.CreateBook(book)
	if err != nil {
		log.Fatal("fail to create book")
	}

	// Create token with `book:create` credential.
	tokenAdmin, err := utils.GenerateNewTokens(
		user.ID.String(),
		[]string{"book:create", "book:update", "book:delete"},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define a structure for specifying input and output data of a single test case.
	test := struct {
		route        string // input route
		method       string // input method
		expectedCode int
	}{
		route:        "/v1/book/" + book.ID.String(),
		method:       "DELETE",
		expectedCode: 204,
	}

	req := httptest.NewRequest(test.method, test.route, nil)
	req.Header.Add("Authorization", "Bearer "+tokenAdmin.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	// Perform the request plain with the AppTest.
	resp, err := AppTest.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to sign in user test")
	}

	defer func() {
		err = db.DeleteUser(user.ID)
		if err != nil {
			log.Fatal("fail to delete user")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
}
