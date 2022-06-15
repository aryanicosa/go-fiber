package user

import (
	"github.com/aryanicosa/go-fiber-rest-api/pkg/user/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

// const USERNAME = "john"
// const PASSWORD = "doe"

// const SecretKey = "TheSecretKey"

// // JwtClaim adds email as a claim to the token
// type JwtClaim struct {
// 	Email string
// 	jwt.StandardClaims
// }

// func MiddlewareBasicAuth(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		username, password, ok := r.BasicAuth()
// 		if !ok {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}

// 		isValid := (username == USERNAME && password == PASSWORD)
// 		if !isValid {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("wrong username/password"))
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func MiddlewareBearerAuth(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		claims := &JwtClaim{}

// 		authToken := r.Header.Get("authorization")
// 		splitToken := strings.Split(authToken, "Bearer ")
// 		authToken = splitToken[1]

// 		token, err := jwt.ParseWithClaims(authToken, claims, func(t *jwt.Token) (interface{}, error) {
// 			return []byte(SecretKey), nil
// 		})

// 		if err != nil {
// 			fmt.Println(err)
// 			if err == jwt.ErrSignatureInvalid {
// 				w.WriteHeader(http.StatusUnauthorized)
// 				w.Write([]byte("unauthorized access"))
// 				return
// 			}

// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write([]byte("wrong auth type"))
// 			return
// 		}

// 		if !token.Valid {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("unauthorized access"))
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	}
// }

type Service struct {
	DB *gorm.DB
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Password string `json:"password"`
}

// GenerateUUID generates new UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func (s *Service) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/user", s.CreateUser)
	// api.Post("/user/:id", s.LoginUser)
	api.Get("/users", s.GetUsers)
	api.Get("/user/:id", s.GetUser)
	api.Put("/user/:id", s.UpdateUser)
	api.Delete("/user/:id", s.DeleteUser)
}

func (s *Service) CreateUser(c *fiber.Ctx) error {
	user := User{}

	err := c.BodyParser(&user)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	validator := validator.New()
	err = validator.Struct(User{})

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	err = s.DB.Create(&model.Users{
		ID:       GenerateUUID(),
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
		Password: password,
	}).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create user"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user created!",
	})
	return nil
}

func (s *Service) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := &model.Users{}

	if id == "" {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "user id can not empty",
		})
		return nil
	}

	err := s.DB.Where("id = ?", id).First(user).Error
	if err != nil {
		c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "user not found",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"data": user,
	})

	return nil
}

func (s *Service) GetUsers(c *fiber.Ctx) error {
	users := &[]model.Users{}

	err := s.DB.Find(users).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get users from db"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"data": users,
	})
	return nil
}

func (s *Service) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userModel := &model.Users{}
	user := User{}

	if id == "" {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "user id can not empty",
		})
		return nil
	}

	err := c.BodyParser(&user)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	validator := validator.New()
	err = validator.Struct(User{})

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	err = s.DB.Where("id = ?", id).First(userModel).Error
	if err != nil {
		c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "user not found",
		})
		return err
	}

	err = s.DB.Updates(&model.Users{
		ID:       userModel.ID,
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
		Password: password,
	}).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update user"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user updated",
	})

	return nil
}

func (s *Service) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := &model.Users{}

	if id == "" {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "user id can not empty",
		})
		return nil
	}

	err := s.DB.Where("id = ?", id).Delete(user).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to delete user",
		})
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user deleted",
	})

	return nil
}

// func (db *Database) LoginUser(context *fiber.Ctx) error {
// 	user := model.User{}

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&user); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	//fmt.Println(user.Email)

// 	var hashedPassword []byte
// 	err := db.QueryRow("SELECT password FROM user WHERE email = ?", user.Email).
// 		Scan(&hashedPassword)
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.Password)); err != nil {
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte("wrong password"))
// 		return
// 	}

// 	claims := &JwtClaim{
// 		Email: user.Email,
// 		StandardClaims: jwt.StandardClaims{
// 			Issuer:    strconv.Itoa(int(user.ID)),
// 			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedToken, err := token.SignedString([]byte(SecretKey))

// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("wcan not login"))
// 		return
// 	}

// 	data := make(map[string]string)

// 	data["message"] = "succesfully login"
// 	data["token"] = signedToken

// 	err = json.NewEncoder(w).Encode(data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
