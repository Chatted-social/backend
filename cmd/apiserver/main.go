package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Chatted-social/backend/handler"
	helper "github.com/Chatted-social/backend/internal/app"
	"github.com/Chatted-social/backend/storage"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var secret = []byte("alskdhjasiudhqwiuhedjkahdkaskdmnknfn")

var (
	port      = flag.String("port", "7070", "which port should be used for webrtc")
	pgURL     = flag.String("pg", os.Getenv("PG_URL"), "url to your postgres database")
	redisAddr = flag.String("redisAddr", os.Getenv("redisAddr"), "ip of redis database")
)

func main() {
	flag.Parse()

	app := helper.NewFiber()
	db, err := storage.Open(*pgURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println(*redisAddr)
	cache, err := storage.NewRedisCache(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()
	store := session.New(session.Config{
		Storage:    cache,
		Expiration: 1 * time.Hour,
	})

	h := handler.NewHandler(handler.Handler{DB: db, RedisCache: cache, SessionsStore: store})

	h.Register(app.Group("/api/auth"), &handler.AuthService{Secret: secret})
	h.Register(app.Group("/api/wall"), &handler.PostStorage{Secret: secret})

	log.Fatal(app.Listen(":" + *port))
}
