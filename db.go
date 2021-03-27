package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	var db *gorm.DB
	var err error

	if conf.Dev {
		db, err = gorm.Open(sqlite.Open("anote.db"), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(conf.PostgreSQL), &gorm.Config{})
	}

	if err != nil {
		log.Printf("[initDb] error: %s", err)
	}

	if err := db.AutoMigrate(&KeyValue{}, &User{}); err != nil {
		panic(err.Error())
	}

	return db
}
