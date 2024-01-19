package database

import (
	"jwtreact/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	dsn := "host=localhost dbname=rda-be port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// connection, err := gorm.Open(mysql.Open("root:@/go_jwt_react"), &gorm.Config{})

	if err != nil {
		panic("Could Not connect to database")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Role{})
	// connection.AutoMigrate(&models.Service{})
	connection.AutoMigrate(&models.UserHasRoles{})

}
