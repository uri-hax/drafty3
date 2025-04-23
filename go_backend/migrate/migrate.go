package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite" 
	"gorm.io/gorm"
)

type ProductGorm struct {
	ProductID          uint      `gorm:"primaryKey;autoIncrement;column:product_id"`
	ProductName        string    `gorm:"type:varchar(255);not null;column:product_name"`
	ProductDescription string    `gorm:"type:text;column:product_description"`
	Price              float64   `gorm:"type:decimal(10,2);not null;column:price"`
	StockQuantity      int       `gorm:"column:stock_quantity;default:0"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}

func (ProductGorm) TableName() string {
	return "products_gorm"
}

func main() {
	db, err := gorm.Open(sqlite.Open("test_migrate.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	if err := db.AutoMigrate(&ProductGorm{}); err != nil {
		log.Fatal("failed to automigrate: ", err)
	}

	log.Println("Automigration complete. Table 'products_gorm' is ready.")
}
