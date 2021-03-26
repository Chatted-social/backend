package main

import (
	"flag"
	"github.com/Chatted-social/backend/handler"
	"github.com/Chatted-social/backend/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/jwt/v2"
	"log"
)

var secret = []byte("alskdhjasiudhqwiuhedjkahdkaskdmnknfn")

var port = flag.String("port", "7070", "which port should be used for server")

//TODO: vernutb proshloe
//flag.String("PG_URL", os.Getenv("PG_URL"), "url to your postgres database")
var PGURL = "host=localhost port=5432 user=admin password=admin dbname=postgres sslmode=disable"

func main() {
	flag.Parse()

	app := newFiber()

	db, err := storage.Open(PGURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	h := handler.NewHandler(handler.Handler{DB: db})

	// example of jwt required
	b := app.Group("/api/blog")
	b.Use(jwtware.New(jwtware.Config{
		SigningKey: secret,
	}))

	h.Register(app.Group("/api/auth"), &handler.AuthService{Secret: secret})
	h.Register(app.Group("/api/wall"), &handler.PostStorage{Secret: secret})

	log.Fatal(app.Listen(":" + *port))
}

func newFiber() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	return app
}
