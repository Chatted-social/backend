package handler

import (
	"errors"
	"github.com/Chatted-social/backend/app"
	j "github.com/Chatted-social/backend/jwt"
	"github.com/Chatted-social/backend/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"strings"
)

type PostStorage struct {
	handler
	Secret []byte
}

func (s PostStorage) REGISTER(h handler, g fiber.Router) {
	s.handler = h
	g.Post("/create", s.Create)
	g.Patch("/update", s.Update)
	g.Delete("/delete/:id", s.Delete)
	g.Get("/post/:id", s.ByID)
	g.Get("/posts/:ids", s.PostsByIDs)
	g.Get("/user/posts/:user_id", s.UserPosts)
}

type CreatePostForm struct {
	Title string `json:"title" validate:"required,min=3,max=32"`
	Body  string `json:"body" validate:"required,min=3,max=255"`
}

type UpdatePatchForm struct {
	ID    int    `sq:"id" json:"id"`
	Title string `json:"title" validate:"required,min=3,max=32"`
	Body  string `json:"body" validate:"required,min=3,max=32"`
}

func (s PostStorage) Create(c *fiber.Ctx) error {
	form := CreatePostForm{}

	if err := c.BodyParser(&form); err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}
	if err := Validate(&form); err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]

	tokenx, err := jwt.ParseWithClaims(token, &j.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(app.Err(err.Error()))
	}

	if !tokenx.Valid {
		return errors.New("handler: token is not valid")
	}

	claims, ok := tokenx.Claims.(*j.Claims)

	if !ok {
		return errors.New("handler: claims ne ok")
	}

	p, err := s.db.Posts.Create(storage.Post{
		Body:    form.Body,
		Title:   form.Title,
		OwnerID: claims.UserID,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(p)
}

func (s PostStorage) Update(c *fiber.Ctx) error {
	form := UpdatePatchForm{}

	if err := c.BodyParser(&form); err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}
	if err := Validate(&form); err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]

	tokenx, err := jwt.ParseWithClaims(token, &j.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(app.Err(err.Error()))
	}

	if !tokenx.Valid {
		return errors.New("handler: token is not valid")
	}

	claims, ok := tokenx.Claims.(*j.Claims)

	if !ok {
		return errors.New("handler: claims ne ok")
	}

	post, err := s.db.Posts.Update(storage.Post{
		OwnerID: claims.UserID,
		Title:   form.Title,
		Body:    form.Body,
		ID:      form.ID,
	})

	if err != nil {
		return c.Status(http.StatusNotAcceptable).JSON(app.Err(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(post)

}

func (s PostStorage) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]

	tokenx, err := jwt.ParseWithClaims(token, &j.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(app.Err(err.Error()))
	}

	if !tokenx.Valid {
		return errors.New("handler: token is not valid")
	}

	claims, ok := tokenx.Claims.(*j.Claims)

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(app.Err("claims is not ok"))
	}

	err = s.db.Posts.Delete(storage.Post{
		OwnerID: claims.UserID,
		ID:      id,
	})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s PostStorage) ByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	post, err := s.db.Posts.ByID(id)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(app.Err(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(post)
}

func (s PostStorage) PostsByIDs(c *fiber.Ctx) error {
	ids := strings.Split(c.Params("ids"), ",")

	var posts []storage.Post

	for _, id := range ids {
		id, err := strconv.Atoi(id)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
		}

		post, err := s.db.Posts.ByID(id)

		if err != nil {
			continue
		}

		posts = append(posts, post)

	}

	return c.Status(http.StatusOK).JSON(posts)
}

func (s PostStorage) UserPosts(c *fiber.Ctx) error {
	user_id, err := strconv.Atoi(c.Params("user_id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	limit, err := strconv.Atoi(c.Query("limit", "100"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	posts, err := s.db.Posts.UserPosts(user_id, limit, offset)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(app.Err(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(posts)

}
