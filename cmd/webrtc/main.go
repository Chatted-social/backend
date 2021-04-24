package main

import (
	"flag"
	"log"
	"os"

	helper "github.com/Chatted-social/backend/internal/app"
	"github.com/Chatted-social/backend/internal/wserver"
	"github.com/Chatted-social/backend/storage"
	"github.com/Chatted-social/backend/webrtc"
)

var port = flag.String("port", "7070", "which port should be used for webrtc")

var PGURL = flag.String("PG_URL", os.Getenv("PG_URL"), "url to your postgres database")

func main() {
	flag.Parse()

	app := helper.NewFiber()

	db, err := storage.Open(*PGURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	sh := webrtc.NewHandler(webrtc.Handler{DB: db})

	s := wserver.NewServer(wserver.Settings{UseJWT: false, OnError: sh.OnError})

	s.Handle(webrtc.EventTypeJoinRoom, sh.OnJoinRoom)
	s.Handle(webrtc.EventTypeAnswer, sh.OnAnswer)
	s.Handle(webrtc.EventTypeOffer, sh.OnOffer)
	s.Handle(webrtc.EventTypeIceCandidate, sh.OnIceCandidate)
	s.Handle(wserver.OnDisconnect, sh.OnDisconnect)

	app.Get("/ws", s.Listen())

	log.Fatal(app.Listen(":" + *port))
}
