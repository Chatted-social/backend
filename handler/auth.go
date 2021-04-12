
package handler

import (
	"database/sql"
	"github.com/Chatted-social/backend/app"
	"github.com/Chatted-social/backend/jwt"
	"github.com/Chatted-social/backend/storage"
	"github.com/Chatted-social/backend/validator"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

type AuthService struct {
	handler
	Secret []byte
}

func (s AuthService) REGISTER(h handler, g fiber.Router) {
	s.handler = h
	g.Post("/login", s.Login)
	g.Post("/register", s.Register)
}

type RegisterForm struct {
	Email     string `json:"email" validate:"required,email,min=6"`
	Username  string `json:"username" validate:"required,min=3,max=32"`
	FirstName string `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string `json:"last_name" validate:"required,min=3,max=32"`
	Password  string `json:"password" validate:"required,min=3,max=32"`
}

func Validate(i interface{}) error {
	validate := validator.New()
	err := validate.Validate(i)
	return err
}

func (s AuthService) Register(c *fiber.Ctx) error {
	form := RegisterForm{}
	if err := c.BodyParser(&form); err != nil {
		return err
	}
	if err := Validate(&form); err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	// removing caps, because we don't need it
	form.Username = strings.ToLower(form.Username)

	exists, err := s.db.Users.ExistsByUsername(form.Username)
	if err != nil {
		return err
	}
	if exists {
		return c.Status(http.StatusConflict).JSON(app.Err("username/email already taken"))
	}
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u := storage.User{
		Email:             form.Email,
		Username:          form.Username,
		FirstName:         form.FirstName,
		LastName:          form.LastName,
		EncryptedPassword: string(encryptedPass),
	}
	err = s.db.Users.Create(u)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(u)
}

func (s AuthService) Login(c *fiber.Ctx) error {
	var form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&form); err != nil {
		return err
	}
	if err := Validate(&form); err != nil {
		return err
	}
	form.Username = strings.ToLower(form.Username)
	user, err := s.db.Users.ByUsername(form.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return c.Status(http.StatusNotFound).JSON(app.Err("user not found"))
	}

	if !s.compareHash(form.Password, user.EncryptedPassword) {
		return c.Status(http.StatusBadRequest).JSON(app.Err("wrong credentials"))
	}

	// todo: implement this
	//if user.BannedUntil.Sub(time.Now()).Seconds() > 0 {
	//	return c.JSON(http.StatusForbidden, app.Err("user is restricted"))
	//}

	token := jwt.NewWithClaims(jwt.Claims{
		UserID: strconv.Itoa(user.ID),
	})

	t, err := token.SignedString(s.Secret)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": t})
}

func (s AuthService) compareHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
