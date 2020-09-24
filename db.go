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

	if err := db.AutoMigrate(&KeyValue{}, &Transaction{}, &User{}, &Shout{}).Error; err != nil {
		panic(err.Error())
	}

	return db
}

// func initDbBak() *gorm.DB {
// 	var db *gorm.DB
// 	var err error

// 	if conf.Dev {
// 		db, err = gorm.Open("sqlite3", "anote.db")
// 	} else {
// 		db, err = gorm.Open("postgres", strings.Replace(conf.PostgreSQL, "anotenew", "anote", 1))
// 	}

// 	if err != nil {
// 		log.Printf("[initDb] error: %s", err)
// 	}

// 	db.LogMode(conf.Debug)

// 	// if err := db.AutoMigrate(&KeyValue{}, &Transaction{}, &User{}, &Shout{}).Error; err != nil {
// 	// 	panic(err.Error())
// 	// }

// 	return db
// }

// func restoreBackup() {
// 	var users []*User
// 	db.Find(&users)
// 	log.Println(len(users))
// 	for i, u := range users {
// 		uB := &User{TelegramID: u.TelegramID}
// 		dbBak.First(uB, uB)
// 		if u.ID != uB.ID {
// 			var refUsers []*User
// 			dbBak.Where(&User{ReferralID: uB.ID}).Find(&refUsers)
// 			for _, ur := range refUsers {
// 				urn := &User{TelegramID: ur.TelegramID}
// 				db.First(urn, urn)
// 				if urn.ReferralID != u.ID && urn.ReferralID != 0 {
// 					log.Printf("Found %d", i)
// 				}
// 			}
// 		}
// 	}
// }
