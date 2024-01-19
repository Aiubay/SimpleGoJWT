package controllers

import (
	"jwtreact/database"
	"jwtreact/models"

	"github.com/gofiber/fiber/v2"
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

func InsertUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Delete(&users)
	return c.JSON(fiber.Map{
		"Messages": "Success",
	})
}
