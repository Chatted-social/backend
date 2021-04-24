package webrtc

import (
	"github.com/Chatted-social/backend/internal/wserver"
	"github.com/Chatted-social/backend/storage"
	"sync"
)

type Clients struct {
	cl map[string]*Client
	sync.RWMutex
}

func NewClients() *Clients {
	return  &Clients{
		cl: make(map[string]*Client),
	}
}

type Client struct {
	*wserver.Conn `json:"-"`
	ID            string `json:"id"`
}

func (r *Clients) Read(key string) *Client {
	r.Lock()
	defer r.Unlock()
	return r.cl[key]
}

func (r *Clients) Write(key string, c *Client) {
	r.Lock()
	defer r.Unlock()
	r.cl[key] = c
}

func (r *Clients) Delete(key string) {
	r.Lock()
	defer r.Unlock()
	delete(r.cl, key)
}


type Room struct {
	items []*Client
	sync.RWMutex
}

func (r *Room) Append(c *Client)  {
	r.Lock()
	defer r.Unlock()
	r.items = append(r.items, c)
}

func (r *Room) Delete(id string)  {
	r.Lock()
	defer r.Unlock()
	for i, cl := range r.items {
		if cl.ID == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
		}
	}
}

type Rooms struct {
	rms map[string]*Room
	sync.RWMutex
}

func (r *Rooms) Read(key string) *Room {
	r.Lock()
	defer r.Unlock()
	return r.rms[key]
}

func (r *Rooms) Write(key string, rm *Room) {
	r.Lock()
	defer r.Unlock()
	r.rms[key] = rm
}

func (r *Rooms) Delete(key string) {
	r.Lock()
	defer r.Unlock()
	delete(r.rms, key)
}

func NewRooms() *Rooms{
	return &Rooms{rms: make(map[string]*Room)}
}


type Handler struct {
	DB     *storage.DB
}

type handler struct {
	db     *storage.DB
	rooms *Rooms
	clients *Clients
}

func NewHandler(h Handler) *handler {
	return &handler{
		db:     h.DB,
		rooms: NewRooms(),
		clients: NewClients(),
	}
}