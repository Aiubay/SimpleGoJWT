package controllers

import (
	"jwtreact/database"
	"jwtreact/models"

	"github.com/bxcodec/faker/v4"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func InsertRole(c *fiber.Ctx) error {
	roles := []models.Role{
		{ID: 1, Role: "Admin", Description: "Admin role", Slug: "admin"},
		{ID: 2, Role: "User", Description: "User role", Slug: "user"},
		{ID: 3, Role: "Guest", Description: "Guest role", Slug: "guest"},
	}

	for _, role := range roles {
		result := database.DB.Create(&role)
		if result.Error != nil {
			return c.JSON(fiber.Map{
				"Messages": result.Error,
			})
		}
	}
	return c.JSON(fiber.Map{
		"Messages": "Success",
	})
}

func DeleteAllUsers(c *fiber.Ctx) error {
	// var users []models.User

	database.DB.Exec("Truncate TABLE users")

	return c.JSON(fiber.Map{
		"Messages": "Success",
	})
}

func CreateUsers(c *fiber.Ctx) error {

	for i := 0; i < 10; i++ {

		password, _ := bcrypt.GenerateFromPassword([]byte("testing"), bcrypt.DefaultCost)

		user := models.User{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: password,
		}

		database.DB.Create(&user)

		AssignRole(int(user.ID), "guest")

	}

	return c.JSON(fiber.Map{
		"Messages": "Success",
	})
}
