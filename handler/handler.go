package handler

import (
	"github.com/Chatted-social/backend/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Handler struct {
	DB            *storage.DB
	Secret        []byte
	RedisCache    *storage.RedisCache
	SessionsStore *session.Store
}

type handler struct {
	db       *storage.DB
	secret   []byte
	cache    *storage.RedisCache
	sessions *session.Store
}

func NewHandler(h Handler) *handler {
	return &handler{
		db:       h.DB,
		secret:   h.Secret,
		cache:    h.RedisCache,
		sessions: h.SessionsStore,
	}
}

type Service interface {
	REGISTER(h handler, g fiber.Router)
}

func (h handler) Register(group fiber.Router, service Service) {
	service.REGISTER(h, group)
}
