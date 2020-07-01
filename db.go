package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDb() *gorm.DB {
	var db *gorm.DB
	var err error

	if conf.Dev {
		db, err = gorm.Open("sqlite3", "anote.db")
	} else {
		db, err = gorm.Open("postgres", conf.PostgreSQL)
	}

	if err != nil {
		log.Printf("[initDb] error: %s", err)
	}

	db.LogMode(conf.Debug)
	db.AutoMigrate(&KeyValue{}, &Transaction{}, &User{})

	return db
}
