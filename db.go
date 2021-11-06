package main

import (
	"fmt"
	"os"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Quiz struct {
	Genre		string `json:"genre"`
	Num			string `json:"num"`
	Caught		string `json:"caught"`
	Flag		string `json:"flag"`
}

func db_connect() *gorm.DB {
	url := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalln("connection failed.", err)
	} else {
		fmt.Println("connected!")
	}

	return db
}
