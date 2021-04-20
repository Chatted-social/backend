package webrtc

import (
	"errors"
	"fmt"
	"github.com/Chatted-social/backend/wserver"
	"github.com/google/uuid"
)

func (h *handler) OnJoinRoom(c *wserver.Context) error {
	var form = EventJoinRoom{}

	if err := c.Bind(&form); err != nil{
		return err
	}

	id := uuid.New().String()

	room := h.rooms.Read(form.RoomID)
	if room == nil{
		room = &Room{items: make([]*Client, 1, 1)}
		h.rooms.Write(form.RoomID, room)
	}

	cl := &Client{Conn: c.Conn, ID: id}

	room.Lock()
	for _, u := range room.items{
		if u != nil {
			u.Conn.WriteJSON(&EventUserJoined{UserID: id})
		}
	}
	room.Unlock()

	room.Append(cl)

	form.OtherUsers = room
	form.UserID = id

	h.clients.Write(id,cl )
	c.Update.Data = form
	return c.Conn.WriteJSON(c.Update)
}

func (h *handler) OnOffer(c *wserver.Context) error {
	var form = EventOffer{}

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
	var form = EventAnswer{}

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