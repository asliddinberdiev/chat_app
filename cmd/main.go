package main

import (
	"log"

	"github.com/asliddinberdiev/chat_app/conf"
	"github.com/asliddinberdiev/chat_app/db"
	"github.com/asliddinberdiev/chat_app/internal/user"
	"github.com/asliddinberdiev/chat_app/router"
)

func main() {
	conf.Load(".")

	dbConn, err := db.NewDatabse(conf.Cfg.Postgres)
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}
	userR := user.NewRepository(dbConn.GetDB())
	userS := user.NewService(userR)
	userH := user.NewHandler(userS)

	router.InitRouter(userH)

	if err := router.Start(conf.Cfg.App.Host + ":" + conf.Cfg.App.Port); err != nil {
		log.Fatal(err)
	}
}
