package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open("lib-Mgmt?_fk=1"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Library{})
	database.AutoMigrate(&Users{})
	database.AutoMigrate(&BookInventory{})
	database.AutoMigrate(&RequestEvents{})
	database.AutoMigrate(&IssueRegistry{})

	DB = database
}
