package config

import (
	"assignment-2/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "root"
	dbPort   = "5432"
	dbName   = "assignment-2"
	db       *gorm.DB
	err      error
)

func StartDB() *gorm.DB {
	config := fmt.Sprintf("host=%s user=%s password=%s dbame=%s port=%s sslmode=disable",
		host, user, password, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		panic("Failed Connect to Database")
	}
	defer fmt.Println("Successfully Connected to Database")

	db.AutoMigrate(models.Order{}, models.Item{})
	return db
}
