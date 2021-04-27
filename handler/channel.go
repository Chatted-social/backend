package handler

import (
	"github.com/Chatted-social/backend/app"
	"github.com/Chatted-social/backend/jwt"
	"github.com/Chatted-social/backend/storage"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"net/http"
	"strconv"
	"strings"
)

type ChannelService struct {
	handler
	Secret []byte
}

func (s ChannelService) REGISTER(h handler, g fiber.Router) {
	s.handler = h

	g.Use(jwtware.New(jwtware.Config{SigningKey: s.Secret}))

	//g.Get("/post/:channel_id/:post_id", s.PostByID)

	//Interaction with channel's posts
	g.Post("/post/create", s.PostCreate)
	//g.Put("/post/update", s.PostUpdate)
	//g.Delete("/post/delete/:id", s.PostDelete)

	//Interaction with channel
	g.Post("/create", s.Create)
	g.Put("/update", s.Update)
	g.Delete("/delete/:id", s.Delete)

	//Channel administration
	g.Post("/promote/:user_id", s.Promote)
	g.Post("/demote/:user_id", s.Demote)
	g.Post("/kick/:user_id", s.Kick)
	g.Post("/block/:user_id", s.Block)

	//Users interaction
	g.Post("/user/subscribe", s.Subscribe)
	g.Post("/user/unsubscribe", s.Unsubscribe)

}

//func (s ChannelService) PostByID(c *fiber.Ctx) error {
//	ChannelID, err := strconv.Atoi(c.Params("channel_id", "0"))
//
//	if err != nil {
//		return err
//	}
//
//	PostID, err := strconv.Atoi(c.Params("post_id", "0"))
//
//	if err != nil {
//		return err
//	}
//
//	s.db.Posts.ByID()
//
//}

//func (s ChannelService) Channel(c *fiber.Ctx) error {
//
//}

func (s ChannelService) PostCreate(c *fiber.Ctx) error {
	var form struct {
		ChannelID int    `json:"channel_id"`
		Body      string `json:"body" validate:"required,min=3,max=255"`
	}
	err := c.BodyParser(&form)
	if err != nil {
		return err
	}
	err = Validate(&form)
	if err != nil {
		return err
	}

	UserID := jwt.From(c.Locals("user")).UserID

	exists, err := s.db.Users.ExistsByID(UserID)

	if err != nil {
		return err
	}

	isOwner, err := s.db.Channels.UserIsOwner(form.ChannelID, UserID)

	if err != nil {
		return err
	}

	isAdmin, err := s.db.Channels.UserIsAdmin(form.ChannelID, UserID)

	if err != nil {
		return err
	}

	if !isOwner && !isAdmin {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	if !exists {
		return c.Status(http.StatusBadRequest).JSON(app.Err("user is not exists"))
	}

	latest := s.db.Channels.LatestPostID(form.ChannelID)

	p, err := s.db.Posts.Create(storage.Post{
		PostID:  latest + 1,
		FromID:  UserID,
		OwnerID: form.ChannelID,
		Body:    form.Body,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(p)

}

//func (s ChannelService) PostUpdate(c *fiber.Ctx) error {
//
//}
//
//func (s ChannelService) PostDelete(c *fiber.Ctx) error {
//
//}

func (s ChannelService) Create(c *fiber.Ctx) error {
	var form struct {
		Title    string `json:"title" validate:"required,min=3,max=35"`
		Username string `json:"username,omitempty" validate:"max=35"`
	}
	err := c.BodyParser(&form)
	if err != nil {
		return err
	}

	err = Validate(&form)
	if err != nil {
		return err
	}

	UserID := jwt.From(c.Locals("user")).UserID

	form.Username = strings.ToLower(form.Username)

	valid := app.UsernameValidator(form.Username)

	if !valid && form.Username != "" {
		return c.Status(http.StatusBadRequest).JSON(app.Err("username is not valid"))
	}

	exists, err := app.UsernameExists(s.db, form.Username)

	if err != nil {
		return err
	}

	if exists {
		return c.Status(http.StatusOK).JSON(app.Err("username already taken"))
	}

	channel, err := s.db.Channels.Create(storage.Channel{
		OwnerID:  UserID,
		Username: form.Username,
		Title:    form.Title,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(channel)

}

func (s ChannelService) Update(c *fiber.Ctx) error {
	var form struct {
		ChannelID int    `json:"channel_id" validate:"required,min=-100000000,max=-1000"`
		Title     string `json:"title" validate:"required,min=3,max=35"`
		Username  string `json:"username" validate:"min=1,max=35"`
	}
	err := c.BodyParser(&form)

	if err != nil {
		return err
	}

	err = Validate(&form)

	if err != nil {
		return err
	}

	UserID := jwt.From(c.Locals("user")).UserID

	err = s.db.Channels.Update(storage.Channel{
		ID:       form.ChannelID,
		Title:    form.Title,
		Username: form.Username,
		OwnerID:  UserID,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s ChannelService) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id", "0"))

	if err != nil {
		return err
	}

	exists, err := s.db.Channels.ExistsByID(id)

	if err != nil {
		return err
	}

	if !exists {
		return c.Status(http.StatusBadRequest).JSON(app.Err("channel is not exists"))
	}

	UserID := jwt.From(c.Locals("user")).UserID

	err = s.db.Channels.Delete(storage.Channel{
		OwnerID: UserID,
		ID:      id,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s ChannelService) Promote(c *fiber.Ctx) error {
	var form struct {
		ChannelID int `json:"channel_id" validate:"required"`
	}

	err := c.BodyParser(&form)
	if err != nil {
		return err
	}
	err = Validate(&form)
	if err != nil {
		return err
	}

	TargetID, err := strconv.Atoi(c.Params("user_id", "0"))

	if err != nil {
		return err
	}

	exists, err := s.db.Users.ExistsByID(TargetID)

	if err != nil {
		return err
	}

	if !exists {
		return c.Status(http.StatusNotFound).JSON(app.Err("user is not exists"))
	}

	UserID := jwt.From(c.Locals("user")).UserID

	if UserID == TargetID {
		return c.Status(http.StatusBadRequest).JSON(app.Err("you have no rights to this action"))
	}

	if exists, _ := s.db.Channels.UserIsSubscriber(form.ChannelID, TargetID); !exists {
		return c.Status(http.StatusBadRequest).JSON(app.Err("user is not a subscriber"))
	}

	if isOwner, _ := s.db.Channels.UserIsOwner(form.ChannelID, UserID); !isOwner {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	if isAdmin, _ := s.db.Channels.UserIsAdmin(form.ChannelID, TargetID); isAdmin {
		return c.Status(http.StatusBadRequest).JSON(app.Err("user is already a admin"))
	}

	err = s.db.Channels.UserPromote(form.ChannelID, TargetID)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s ChannelService) Demote(c *fiber.Ctx) error {
	var form struct {
		ChannelID int `json:"channel_id" validate:"required"`
	}

	err := c.BodyParser(&form)
	if err != nil {
		return err
	}
	err = Validate(&form)
	if err != nil {
		return err
	}

	TargetID, err := strconv.Atoi(c.Params("user_id", "0"))

	if err != nil {
		return err
	}

	exists, err := s.db.Users.ExistsByID(TargetID)

	if err != nil {
		return err
	}

	if !exists {
		return c.Status(http.StatusNotFound).JSON(app.Err("user is not exists"))
	}

	UserID := jwt.From(c.Locals("user")).UserID

	if UserID == TargetID {
		return c.Status(http.StatusBadRequest).JSON(app.Err("you have no rights to this action"))
	}

	if exists, _ := s.db.Channels.UserIsSubscriber(form.ChannelID, TargetID); !exists {
		return c.Status(http.StatusBadRequest).JSON(app.Err("user is not a subscriber"))
	}

	if isOwner, _ := s.db.Channels.UserIsOwner(form.ChannelID, UserID); !isOwner {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	if isAdmin, _ := s.db.Channels.UserIsAdmin(form.ChannelID, TargetID); !isAdmin {
		return c.Status(http.StatusBadRequest).JSON(app.Err("user is not a admin"))
	}

	err = s.db.Channels.UserDemote(form.ChannelID, TargetID)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s ChannelService) Kick(c *fiber.Ctx) error {
	var form struct {
		ChannelID int `json:"channel_id"`
	}

	TargetID, err := strconv.Atoi(c.Params("user_id", "0"))

	if err != nil {
		return err
	}

	UserID := jwt.From(c.Locals("user")).UserID

	isOwner, err := s.db.Channels.UserIsOwner(form.ChannelID, UserID)

	if err != nil {
		return err
	}

	isAdmin, err := s.db.Channels.UserIsAdmin(form.ChannelID, UserID)
	if err != nil {
		return err
	}

	if !isAdmin && !isOwner {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	if isAdmin, _ := s.db.Channels.UserIsAdmin(form.ChannelID, TargetID); isAdmin {
		err := s.db.Channels.UserDemote(form.ChannelID, TargetID)

		if err != nil {
			return err
		}

	}

	if isOwner, _ := s.db.Channels.UserIsOwner(form.ChannelID, TargetID); isOwner {
		return c.Status(http.StatusBadRequest).JSON(app.Err("you have no rights to this action"))
	}

	err = s.db.Channels.UserUnsubscribe(form.ChannelID, TargetID)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s ChannelService) Block(c *fiber.Ctx) error {
	var form struct {
		ChannelID int `json:"channel_id"`
	}
	err := c.BodyParser(&form)
	if err != nil {
		return err
	}
	err = Validate(&form)
	if err != nil {
		return err
	}

	TargetID, err := strconv.Atoi(c.Params("user_id", "0"))

	if err != nil {
		return err
	}

	exists, err := s.db.Users.ExistsByID(TargetID)

	if err != nil {
		return err
	}

	if !exists {
		return c.Status(http.StatusNotFound).JSON(app.Err("user is not exists"))
	}

	UserID := jwt.From(c.Locals("user")).UserID

	if UserID == TargetID {
		return c.Status(http.StatusBadRequest).JSON(app.Err("you have no rights to this action"))
	}

	isAdmin, err := s.db.Channels.UserIsAdmin(form.ChannelID, UserID)
	if err != nil {
		return err
	}

	isOwner, err := s.db.Channels.UserIsOwner(form.ChannelID, UserID)
	if err != nil {
		return err
	}

	if !isOwner && !isAdmin {
		return c.Status(http.StatusForbidden).JSON(app.Err("you have no rights to this action"))
	}

	TargetIsAdmin, err := s.db.Channels.UserIsAdmin(form.ChannelID, TargetID)
	if err != nil {
		return err
	}

	TargetIsOwner, err := s.db.Channels.UserIsOwner(form.ChannelID, TargetID)
	if err != nil {
		return err
	}

	if isAdmin && (TargetIsAdmin || TargetIsOwner) {
		return c.Status(http.StatusOK).JSON(app.Err("you have no rights to this action"))
	}

	if isBlocked := s.db.Channels.UserIsBlocked(form.ChannelID, TargetID); isBlocked {
		return c.Status(http.StatusBadRequest).JSON(app.Err("user is already blocked"))
	}

	if isSubscribed, _ := s.db.Channels.UserIsSubscriber(form.ChannelID, TargetID); isSubscribed {
		err := s.db.Channels.UserUnsubscribe(form.ChannelID, TargetID)

		if err != nil {
			return err
		}

	}

	err = s.db.Channels.UserBlock(form.ChannelID, TargetID)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s ChannelService) Subscribe(c *fiber.Ctx) error {
	var form struct {
		InvitationLink string `json:"invitation_link,omitempty"`
		Username       string `json:"username,omitempty"`
	}

	err := c.BodyParser(&form)
	if err != nil {
		return err
	}

	UserID := jwt.From(c.Locals("user")).UserID

	SubType := c.Query("type", "username")

	if SubType == "username" {
		if exists, _ := s.db.Channels.ExistsByUsername(form.Username); !exists {
			return c.Status(http.StatusNotFound).JSON(app.Err("channel with this username does not exist"))
		}

		Channel, err := s.db.Channels.ByUsername(form.Username)

		if err != nil {
			return err
		}

		if isBlocked := s.db.Channels.UserIsBlocked(Channel.ID, UserID); isBlocked {
			return c.Status(http.StatusForbidden).JSON(app.Err("you are blocked from this channel"))
		}

		err = s.db.Channels.UserSubscribe(Channel.ID, UserID)

		if err != nil {
			return err
		}

	} else if SubType == "link" {
		//TODO: implement this
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}

func (s ChannelService) Unsubscribe(c *fiber.Ctx) error {
	var form struct {
		ChannelID int `json:"channel_id" validate:"required"`
	}

	err := c.BodyParser(&form)
	if err != nil {
		return err
	}
	err = Validate(&form)
	if err != nil {
		return err
	}

	UserID := jwt.From(c.Locals("user")).UserID

	err = s.db.Channels.UserUnsubscribe(form.ChannelID, UserID)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(app.Ok())

}
