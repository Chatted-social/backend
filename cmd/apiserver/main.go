package main

import (
	"flag"
	"github.com/Chatted-social/backend/handler"
	"github.com/Chatted-social/backend/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var secret = []byte("alskdhjasiudhqwiuhedjkahdkaskdmnknfn")

var port = flag.String("port", "7070", "which port should be used for server")

var PGURL = flag.String("PG_URL", os.Getenv("PG_URL"), "url to your postgres database")

var debug = flag.Bool("debug", false, "debug mode")

func init() {
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	flag.Parse()

	app := newFiber()

	db, err := storage.Open(*PGURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	h := handler.NewHandler(handler.Handler{DB: db})

	app.Use(func(c *fiber.Ctx) error {
		log.WithFields(log.Fields{
			"method":      c.Method(),
			"ip_address":  c.IP(),
			"status_code": c.Response().StatusCode(),
			"path":        c.Path(),
		}).Info()

		return c.Next()

	})

	v1 := app.Group("/api/v1")
	h.Register(v1.Group("/auth"), &handler.AuthService{Secret: secret})
	h.Register(v1.Group("/wall"), &handler.PostStorage{Secret: secret})
	h.Register(v1.Group("/channel"), &handler.ChannelService{Secret: secret})
	
	log.Fatal(app.Listen(":" + *port))
}

func newFiber() *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		log.WithFields(log.Fields{
			"ip_address":    c.IP(),
			"path":          c.Path(),
			"method":        c.Method(),
			"error_message": err.Error(),
		}).Debug()

		return c.Status(http.StatusInternalServerError).SendString(err.Error())

	}})

	app.Use(cors.New())

	return app
}
