package database

import (
	"fmt"
	"jwtreact/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s dbname=%s port=%s sslmode=%s timezone=%s", os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_TIMEZONE"))
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	DB = connection
}

func Migrate() {

	if !DB.Migrator().HasTable(&models.User{}) {
		DB.AutoMigrate(&models.User{})
	}
	if !DB.Migrator().HasTable(&models.Role{}) {
		DB.AutoMigrate(&models.Role{})
	}
	// connection.AutoMigrate(&models.Service{})
	if !DB.Migrator().HasTable(&models.UserHasRoles{}) {
		DB.AutoMigrate(&models.UserHasRoles{})
	}
}
