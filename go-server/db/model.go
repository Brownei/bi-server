package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string
	Age   uint16
	Likes string
}

type Message struct {
	gorm.Model
	ID       string
	Msg      string
	Author   string
	Receiver string
}
