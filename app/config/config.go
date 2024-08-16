package config

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/syrlramadhan/pendaftaran-coc/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Inisialisasi database
func InitDB() {
	var err error

	db, err := gorm.Open(sqlite.Open("pendaftar.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = db.AutoMigrate(&models.Pendaftar{})
	if err != nil {
		panic("Error migrating database: " + err.Error())
	}
	fmt.Println("Database connected and migrated successfully")

	DB = db
}