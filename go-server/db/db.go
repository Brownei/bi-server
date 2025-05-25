package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDb(chanDB chan *gorm.DB) {
	dsn := "host=localhost user=tcp_db password=tcp_db dbname=tcp_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to postgres")
	} else {
		log.Print("Connected to database successfully")
	}

	chanDB <- db
}
