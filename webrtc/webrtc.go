package webrtc

import (
	"errors"
	"fmt"
	"github.com/Chatted-social/backend/internal/wserver"
	"github.com/google/uuid"
)

func (h *handler) OnJoinRoom(c *wserver.Context) error {
	var form = EventJoinRoom{}

	if err := c.Bind(&form); err != nil{
		return err
	}

	id, ok := c.Conn.Get("id").(string)
	if ok {
		if id != "" {
			return errors.New("webrtc: client already exists")
		}
	}

	id = uuid.New().String()
	c.Conn.Set("id", id)

	room := h.rooms.Read(form.RoomID)
	if room == nil{
		room = &Room{items: make([]*Client, 0, 2)}
		h.rooms.Write(form.RoomID, room)
	}

	cl := &Client{Conn: c.Conn, ID: id}

	room.Lock()
	for _, u := range room.items{
		if u != nil {
			u.Conn.WriteJSON(&wserver.Update{ EventType: EventTypeUserJoined, Data: &EventUserJoined{UserID: id}})
		}
	}
	room.Unlock()

	room.Append(cl)
	c.Conn.Set("room", form.RoomID)

	form.OtherUsers = room.items
	form.UserID = id

	h.clients.Write(id,cl )
	c.Update.Data = form
	return c.Conn.WriteJSON(c.Update)
}

func (h *handler) OnOffer(c *wserver.Context) error {
	var form = EventHandshake{}

	if err := c.Bind(&form); err != nil{
		return err
	}

	target := h.clients.Read(form.Target)
	if target == nil{
		return errors.New(fmt.Sprintf("webrtc: client with id %s does not exists", form.Target))
	}
	return target.WriteJSON(c.Update)
}

func (h *handler) OnAnswer(c *wserver.Context) error {
	var form = EventHandshake{}

	if err := c.Bind(&form); err != nil{
		return err
	}

	target := h.clients.Read(form.Target)
	if target == nil{
		return errors.New(fmt.Sprintf("webrtc: client with id %s does not exists", form.Target))
	}

	return target.WriteJSON(c.Update)
}

func (h *handler) OnIceCandidate(c *wserver.Context) error {
	var form = EventIceCandidate{}
	if err := c.Bind(&form); err != nil{
		return err
	}
	target := h.clients.Read(form.Target)
	if target == nil{
		return errors.New(fmt.Sprintf("webrtc: client with id %s does not exists", form.Target))
	}

	return target.WriteJSON(c.Update)
}