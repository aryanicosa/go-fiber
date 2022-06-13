package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
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
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Password string `json:"password" validate:"required"`
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

	err = s.DB.Create(User{
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
		Password: string(password),
	}).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create user"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been successfully added",
	})
	return nil
}

func (s *Service) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/user", s.CreateUser)
	// api.Post("/user/:id", s.LoginUser)
	// api.Get("/user", s.GetUsers)
	// api.Get("/user/:id", s.GetUser)
	// api.Put("/user/:id", s.UpdateUser)
	// api.Delete("/user/:id", s.DeleteUser)
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

// func GetUsers(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		users, err := db.Query("SELECT id, name, email, address FROM user")
// 		if err != nil {
// 			panic(err)
// 		}

// 		data := []model.User{}

// 		for users.Next() {
// 			var user model.User

// 			err = users.Scan(&user.ID, &user.Name, &user.Email, &user.Address)
// 			if err != nil {
// 				panic(err)
// 			}

// 			data = append(data, user)
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")

// 		err = json.NewEncoder(w).Encode(data)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

// func GetUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		data := model.User{}

// 		vars := mux.Vars(r)
// 		id, err := strconv.Atoi(vars["userId"])

// 		if err != nil || id < 1 {
// 			http.NotFound(w, r)
// 			return
// 		}

// 		err = db.QueryRow("SELECT id, name, email, address FROM user WHERE id = ?", id).
// 			Scan(&data.ID, &data.Name, &data.Email, &data.Address)
// 		if err != nil {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")

// 		err = json.NewEncoder(w).Encode(data)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

// func UpdateUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		data := model.User{}

// 		vars := mux.Vars(r)
// 		id, err := strconv.Atoi(vars["userId"])

// 		if err != nil || id < 1 {
// 			http.NotFound(w, r)
// 			return
// 		}

// 		decoder := json.NewDecoder(r.Body)
// 		if err := decoder.Decode(&data); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		_, err = db.Exec("UPDATE user SET name = ?, email = ?, address = ? WHERE id = ?",
// 			data.Name, data.Email, data.Address, id)

// 		if err != nil {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}

// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }

// func DeleteUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		data := model.User{}

// 		vars := mux.Vars(r)
// 		id, err := strconv.Atoi(vars["userId"])

// 		if err != nil || id < 1 {
// 			http.NotFound(w, r)
// 			return
// 		}

// 		err = db.QueryRow("SELECT id, name, email, address FROM user WHERE id = ?", id).
// 			Scan(&data.ID, &data.Name, &data.Email, &data.Address)
// 		if err != nil {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}

// 		_, err = db.Exec("DELETE FROM user WHERE id = ?", id)
// 		if err != nil {
// 			http.Error(w, "Failed to delete", http.StatusNotFound)
// 			return
// 		}

// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }
