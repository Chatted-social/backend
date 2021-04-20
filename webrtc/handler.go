package webrtc

import (
	"github.com/Chatted-social/backend/storage"
	"github.com/Chatted-social/backend/wserver"
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
	ID string `json:"id"`
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

type Rooms struct {
	rms map[string]*Room
	sync.Mutex
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
	Secret []byte
}

type handler struct {
	db     *storage.DB
	secret []byte
	rooms *Rooms
	clients *Clients
}

func NewHandler(h Handler) *handler {
	return &handler{
		db:     h.DB,
		secret: h.Secret,
		rooms: NewRooms(),
		clients: NewClients(),
	}
}