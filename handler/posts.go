package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Chatted-social/backend/internal/app"
	j "github.com/Chatted-social/backend/jwt"
	"github.com/Chatted-social/backend/storage"
	"github.com/gofiber/fiber/v2"
)

type PostStorage struct {
	handler
	Secret []byte
}

func (s PostStorage) REGISTER(h handler, g fiber.Router) {
	s.handler = h

	g.Get("/post/:id", s.ByID)
	g.Get("/posts/:ids", s.PostsByIDs)
	g.Get("/user/posts/:user_id", s.UserPosts)

	//g.Use(jwtware.New(jwtware.Config{SigningKey: s.Secret}))

	g.Post("/create", s.Create)
	g.Put("/update", s.Update)
	g.Delete("/delete/:id", s.Delete)

}

type CreatePostForm struct {
	Title string `json:"title" validate:"required,min=3,max=32"`
	Body  string `json:"body" validate:"required,min=3,max=255"`
}

type UpdatePutForm struct {
	ID    int    `sq:"id" json:"id"`
	Title string `json:"title" validate:"required,min=3,max=32"`
	Body  string `json:"body" validate:"required,min=3,max=32"`
}

func (s PostStorage) Create(c *fiber.Ctx) error {
	form := CreatePostForm{}

	if err := c.BodyParser(&form); err != nil {
		return err
	}
	if err := Validate(&form); err != nil {
		return err
	}

	sess, err := s.sessions.Get(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(app.Err("authorization required"))
	}

	userID := sess.Get("id").(int)

	err = sess.Save()
	if err != nil {
		return err
	}

	p, err := s.db.Posts.Create(storage.Post{
		Body:    form.Body,
		Title:   form.Title,
		OwnerID: userID,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(p)
}

func (s PostStorage) Update(c *fiber.Ctx) error {
	form := UpdatePutForm{}

	if err := c.BodyParser(&form); err != nil {
		return err
	}
	if err := Validate(&form); err != nil {
		return err
	}

	userID := j.From(c.Locals("user")).UserID

	post, err := s.db.Posts.Update(storage.Post{
		OwnerID: userID,
		Title:   form.Title,
		Body:    form.Body,
		ID:      form.ID,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(post)

}

func (s PostStorage) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	userID := j.From(c.Locals("user")).UserID

	err = s.db.Posts.Delete(storage.Post{
		OwnerID: userID,
		ID:      id,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s PostStorage) ByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	post, err := s.db.Posts.ByID(id)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(post)
}

func (s PostStorage) PostsByIDs(c *fiber.Ctx) error {
	ids := app.StringSliceToInt(strings.Split(c.Params("ids"), ","))

	posts, err := s.db.Posts.PostsIn(ids)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(posts)
}

func (s PostStorage) UserPosts(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("user_id"))

	if err != nil {
		return err
	}

	limit, err := strconv.Atoi(c.Query("limit", "100"))

	if err != nil {
		return err
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))

	if err != nil {
		return err
	}

	posts, err := s.db.Posts.UserPosts(userId, limit, offset)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(posts)

}
