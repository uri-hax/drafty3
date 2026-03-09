package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"drafty3/go_migration/user_model"
)

func main() {
	dsn := "../../db/users_gorm.db?_pragma=foreign_keys(1)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&user_model.Session{},
		&user_model.Profile{},
	); err != nil {
		log.Fatalf("automigrate: %v", err)
	}

	log.Println("AutoMigrate complete. SQLite database created: users_gorm.db")
}
