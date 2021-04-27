package handler

import (
	"github.com/Chatted-social/backend/app"
	j "github.com/Chatted-social/backend/jwt"
	"github.com/Chatted-social/backend/storage"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type PostStorage struct {
	handler
	Secret []byte
}

func (s PostStorage) REGISTER(h handler, g fiber.Router) {
	s.handler = h

	g.Get("/get/:owner_id_:post_id", s.Post)
	g.Get("/user/posts/:user_id", s.UserPosts)

	g.Use(jwtware.New(jwtware.Config{SigningKey: s.Secret}))

	g.Post("/create", s.Create)
	g.Put("/update", s.Update)
	g.Delete("/delete/:id", s.Delete)

	g.Put("/like", s.Like)
	//g.Put("/comment", s.Comment)

}

type CreatePostForm struct {
	OwnerID int    `json:"owner_id" validate:"required"`
	Body    string `json:"body" validate:"required,min=3,max=255"`
}

type UpdatePutForm struct {
	OwnerID int    `json:"owner_id" validate:"required"`
	PostID  int    `json:"post_id" validate:"required"`
	Body    string `json:"body" validate:"required,min=3,max=32"`
}

func (s PostStorage) Create(c *fiber.Ctx) error {
	form := CreatePostForm{}

	if err := c.BodyParser(&form); err != nil {
		return err
	}
	if err := Validate(&form); err != nil {
		return err
	}

	UserID := j.From(c.Locals("user")).UserID

	if exists, _ := s.db.Users.ExistsByID(UserID); !exists {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	latest := s.db.Users.LatestPostID(form.OwnerID)

	p, err := s.db.Posts.Create(storage.Post{
		PostID:  latest + 1,
		FromID:  UserID,
		OwnerID: form.OwnerID,
		Body:    form.Body,
	})

	if err != nil {
		return err
	}

	log.Debugf("created new post %v", p)
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

	UserID := j.From(c.Locals("user")).UserID

	if exists, _ := s.db.Users.ExistsByID(UserID); !exists {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	post, err := s.db.Posts.Update(storage.Post{
		OwnerID: form.OwnerID,
		FromID:  UserID,
		Body:    form.Body,
		PostID:  form.PostID,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(post)

}

func (s PostStorage) Delete(c *fiber.Ctx) error {

	var form struct {
		OwnerID int `json:"owner_id" validate:"required"`
	}

	err := c.BodyParser(&form)
	if err != nil {
		return err
	}
	err = Validate(&form)
	if err != nil {
		return err
	}

	PostID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	UserID := j.From(c.Locals("user")).UserID

	if exists, _ := s.db.Users.ExistsByID(UserID); !exists {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	err = s.db.Posts.Delete(storage.Post{
		OwnerID: form.OwnerID,
		FromID:  UserID,
		PostID:  PostID,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s PostStorage) Like(c *fiber.Ctx) error {
	var form struct {
		OwnerID int `json:"owner_id" validate:"required"`
		PostID  int `json:"post_id" validate:"required"`
	}

	err := c.BodyParser(&form)
	if err != nil {
		return err
	}
	err = Validate(&form)
	if err != nil {
		return err
	}

	UserID := j.From(c.Locals("user")).UserID

	if exists, _ := s.db.Users.ExistsByID(UserID); !exists {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	p := storage.Post{
		OwnerID: form.OwnerID,
		PostID:  form.PostID,
	}

	if liked, _ := s.db.Posts.Liked(UserID, p); liked {
		err := s.db.Posts.Unlike(UserID, p)

		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(app.Ok())

	}

	err = s.db.Posts.Like(UserID, p)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

//TODO: implement this
//func (s PostStorage) Comment(c *fiber.Ctx) error {
//}

//TODO: implement this
//func (s PostStorage) PostsByIDs(c *fiber.Ctx) error {
//	IDs := app.StringSliceToInt(strings.Split(c.Params("ids"), ","))
//
//	posts, err := s.db.Posts.PostsIn(IDs)
//
//	if err != nil {
//		return err
//	}
//
//	return c.Status(http.StatusOK).JSON(posts)
//}

func (s PostStorage) Post(c *fiber.Ctx) error {
	OwnerID, err := strconv.Atoi(c.Params("owner_id"))

	if err != nil {
		return err
	}

	PostID, err := strconv.Atoi(c.Params("post_id"))

	if err != nil {
		return err
	}

	p, err := s.db.Posts.ByIDs(OwnerID, PostID)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(p)

}

func (s PostStorage) UserPosts(c *fiber.Ctx) error {
	UserID, err := strconv.Atoi(c.Params("user_id"))

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

	posts, err := s.db.Posts.UserPosts(UserID, limit, offset)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(posts)

}
