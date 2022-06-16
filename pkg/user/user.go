package user

import (
	"net/http"
	"os"
	"time"

	"github.com/aryanicosa/go-fiber-rest-api/pkg/user/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

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

func (s *Service) LoginUser(c *fiber.Ctx) error {
	user := &model.Users{}

	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	err := c.BodyParser(&input)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	validator := validator.New()
	err = validator.Struct(LoginInput{})

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	err = s.DB.Where("email = ?", input.Email).First(user).Error
	if err != nil {
		c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "user not found",
		})
		return err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	claims := jwt.MapClaims{
		"email": input.Email,
		"iss":   user.ID.String(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not login",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"access_token": signedToken,
	})
}

func (s *Service) LogoutUser(c *fiber.Ctx) error {
	// TODO : find a way to revoke jwt token
	
	// headers := c.GetReqHeaders()
	// tokenRead := strings.Replace(headers["Authorization"], "Bearer ", "", -1)

	// claims := jwt.MapClaims{}
	// _, err := jwt.ParseWithClaims(tokenRead, claims, func(t *jwt.Token) (interface{}, error) {
	// 	return []byte(os.Getenv("SECRET_KEY")), nil
	// })
	// if err != nil {
	// 	c.Status(http.StatusUnprocessableEntity).JSON(
	// 		&fiber.Map{"message": "request failed"})
	// 	return err
	// }

	// for key, val := range claims {
	// 	fmt.Printf("Key: %v, value: %v\n", key, val)
	// }

	// claims["exp"] := time.Now().Unix()
	// fmt.Println(claims)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "already log out",
	})
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
