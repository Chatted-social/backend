// tiny library for easy websocket usage
package wserver

import (
	"encoding/json"
	j "github.com/Chatted-social/backend/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"sync"
	"time"
)

const (
	OnOther = iota
	OnConnect
	OnDisconnect
)

var PongTimeout = time.Second * 60

type HandlerFunc func(ctx *Context) error

type OnErrorFunc func(error, *Context)

type Update struct {
	EventType string      `json:"event_type" mapstructure:"event_type"`
	Data      interface{} `json:"data" mapstructure:"data"`
}


type Settings struct {
	// secret for jwt, pass an empty string if u don't use it
	Secret []byte

	// if true, Context.Get("token") will return token
	UseJWT  bool

	// All errors from handlers goes here
	OnError OnErrorFunc

	//
	Claims  jwt.Claims
}

type Server struct {
	OnError  OnErrorFunc
	handlers map[interface{}]HandlerFunc
	useJWT   bool
	secret   []byte
	claims   jwt.Claims
}

// Conn is websocket.Conn wrapper with mutex
// Do not use Ws without locking mutex, it can
// cause concurrent write issue
type Conn struct {
	Ws *websocket.Conn
	sync.Mutex
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{Ws: conn}
}

func (c *Conn) WriteMessage(code int, msg []byte) error {
	c.Lock()
	defer c.Unlock()
	return c.Ws.WriteMessage(code, msg)
}

func (c *Conn) WriteJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.Ws.WriteJSON(v)
}

func (c *Conn) SetPongHandler(h func(appdata string) error) {
	c.Lock()
	defer c.Unlock()
	c.Ws.SetPongHandler(h)
}

func (c *Conn) Close() error {
	c.Lock()
	defer c.Unlock()
	return c.Close()
}

func NewServer(s Settings) *Server {
	if s.UseJWT && len(s.Secret) < 1 {
		panic("wserver: secret can not be empty string if UseJWT enabled")
	}
	return &Server{
		useJWT:   s.UseJWT,
		OnError:  s.OnError,
		handlers: make(map[interface{}]HandlerFunc),
		secret:   s.Secret,
		claims:   s.Claims,
	}
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// Handle registers a new handler
// If the client sends eventType, Server will call HandlerFunc with specific eventType
func (s *Server) Handle(eventType interface{}, h HandlerFunc, middleware ...MiddlewareFunc) {
	switch eventType.(type) {
	case int:
		break
	case string:
		break
	default:
		panic("wserver: unsupported eventType")
	}
	h = applyMiddleware(h, middleware...)
	s.handlers[eventType] = h
}

// runHandler runs HandlerFunc h with Context c
func (s *Server) runHandler(h HandlerFunc, c *Context) {
	f := func() {
		if err := h(c); err != nil {
			if s.OnError != nil {
				s.OnError(err, c)
			} else {
				log.Println(err)
			}
		}
	}
	f()
}

// Listen is handler that upgrades http client to websocket client
func (s *Server) Listen() fiber.Handler  {
	return websocket.New(func(c *websocket.Conn) {
		var token string
		if s.useJWT {
			token = c.Params("token")
			if token == "" {
				return
			}
			t, err := jwt.ParseWithClaims(token, &j.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return s.secret, nil
			})
			if err != nil {
				return
			}
			if !t.Valid {
				return
			}

		}
			conn := NewConn(c)
			//go s.keepAlive(conn, PongTimeout)
			s.reader(conn, token)
	})
}

// keepAlive will write to client PingMessage to make sure that client is alive.
// Required by websockets documentation.
func (s *Server) keepAlive(conn *Conn, timeout time.Duration) {
	lastResponse := time.Now()
	conn.SetPongHandler(func(_ string) error {
		lastResponse = time.Now()
		return nil
	})
	for {
		err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
		if err != nil {
			conn.Close()
			return
		}
		time.Sleep((timeout * 9) / 10)
		if time.Now().Sub(lastResponse) > timeout {
			log.Printf("Ping don't get response, disconnecting to %s", conn.Ws.LocalAddr())
			err = conn.Close()
			if s.OnError != nil {
				s.OnError(err, nil)
			}
			return
		}
	}
}

// reader Listens to websocket messages, creates context and activates handlers
func (s *Server) reader(conn *Conn, token string) {
	ctx := &Context{Conn: conn, storage: make(map[string]interface{})}
	ctx.Set("token", token)
	s.runOnConnectHandler(ctx)
	for {
		ctx := &Context{Conn: conn, storage: make(map[string]interface{})}
		ctx.Set("token", token)
		_, msg, err := conn.Ws.ReadMessage()
		if err != nil {
			s.OnError(err, ctx)
			s.runOnDisconnectHandler(ctx)
			return
		}

		s.processUpdate(msg, ctx)
	}
}

func (s Server) runOnDisconnectHandler(ctx *Context) {
	h, ok := s.handlers[OnDisconnect]
	if ok {
		s.runHandler(h, ctx)
	}
}

func (s *Server) runOnConnectHandler(ctx *Context) {
	h, ok := s.handlers[OnConnect]
	if ok {
		s.runHandler(h, ctx)
	}
}

// processUpdate converting msg to Update and puts it into Context
// then runs handler with eventType that equals Update.EventType
func (s *Server) processUpdate(msg []byte, c *Context) {
	u := &Update{}
	err := json.Unmarshal(msg, u)
	if err != nil {
		if s.OnError != nil {
			s.OnError(err, c)
		}
	}
	c.Update = u

	handler, ok := s.handlers[u.EventType]
	if !ok {
		h, ok := s.handlers[OnOther]
		if ok {
			handler = h
		} else {
			return
		}
	}
	if handler != nil {
		s.runHandler(handler, c)
	}
}
