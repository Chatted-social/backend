package server

import "github.com/Chatted-social/backend/storage"

type Handler struct {
	DB     *storage.DB
	Secret []byte
}

type handler struct {
	db     *storage.DB
	secret []byte
}

func NewHandler(h Handler) *handler {
	return &handler{
		db:     h.DB,
		secret: h.Secret,
	}
}