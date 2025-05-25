package main

import (
	"log"

	"github.com/brownei/fast-server/cmd"
	"github.com/brownei/fast-server/db"
	"gorm.io/gorm"
)

func main() {
	chanDB := make(chan *gorm.DB)
	go db.InitializeDb(chanDB)

	server := cmd.Server{
		DB:      <-chanDB,
		Clients: make(map[string]*cmd.Client),
	}

	server.Run()
	log.Print("Listening already in 3000")
}
