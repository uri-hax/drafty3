package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"drafty3/go_migration/user_model"
)

// function to create sqlite db for users db from gorm automigrate and log any errors
func main() {
	// where to make new sqlite db - can be changed as needed
	dsn := "../../db/users_gorm.db?_pragma=foreign_keys(1)"
	// open new sqlite db using gorm
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}

	// automigrate all models
	if err := db.AutoMigrate(
		&user_model.Session{},
		&user_model.Profile{},
	); err != nil {
		// log any errors
		log.Fatalf("automigrate: %v", err)
	}

	// log success
	log.Println("AutoMigrate complete. SQLite database created: users_gorm.db")
}
