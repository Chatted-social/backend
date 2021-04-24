package webrtc

import "github.com/Chatted-social/backend/internal/wserver"

func (h handler) OnDisconnect(c *wserver.Context) error  {
	id, ok := c.Conn.Get("id").(string)
	if ok {
		h.clients.Delete(id)
	}

	roomID, ok := c.Conn.Get("room").(string)
	if ok {
		room := h.rooms.Read(roomID)
		if room != nil {
			room.Delete(id)
		}
	}

	return nil
}