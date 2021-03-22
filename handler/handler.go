package handler

import (
	"github.com/Chatted-social/backend/storage"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	DB      *storage.DB
	Secret  []byte
}

type handler struct {
	db *storage.DB
	secret  []byte
}

func NewHandler(h Handler) *handler {
	return &handler{
		db: h.DB,
		secret:  h.Secret,
	}
}

type Service interface {
	REGISTER(h handler, g fiber.Router)
}

func (h handler) Register(group fiber.Router, service Service) {
	service.REGISTER(h, group)
}
