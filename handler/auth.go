package handler

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Chatted-social/backend/internal/app"
	"github.com/Chatted-social/backend/storage"
	"github.com/Chatted-social/backend/validator"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
	_, err = s.db.Users.Create(u)
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

	sess, err := s.createSession(map[string]interface{}{
		"email":      user.Email,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"id":         user.ID,
	})
	if err != nil {
		return err
	}

	sessID := uuid.New().String()
	err = s.cache.Set(sessID, sess, s.sessions.Expiration)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{Name: s.sessions.CookieName,
		Value:   sessID,
		Expires: time.Now().Add(time.Hour),
		Path:    "/"})
	return c.Status(http.StatusOK).JSON(app.Ok())
}

func (s AuthService) compareHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (s AuthService) createSession(info map[string]interface{}) ([]byte, error) {
	var buff = new(bytes.Buffer)

	enc := gob.NewEncoder(buff)
	err := enc.Encode(info)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
