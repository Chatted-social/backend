package webrtc

import (
	"github.com/Chatted-social/backend/wserver"
	"log"
)

func (h *handler) OnError(err error, c *wserver.Context) {
	log.Println(err, c)
	return
}