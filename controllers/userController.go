package controllers

import (
	"jwtreact/database"
	"jwtreact/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func Users(c *fiber.Ctx) error {
	var user []models.User

	result := database.DB.Preload(clause.Associations).Find(&user)

	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": result.Error,
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"message": "No users found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "OK",
		"code":   fiber.StatusOK,
		"data":   user,
	})
}

func GetUserByID(c *fiber.Ctx) error {

	var user models.User
	var hasRoles models.UserHasRoles
	database.DB.Model(&user).Association(clause.Associations).Find(&hasRoles)

	return c.JSON(fiber.Map{
		"data":  user,
		"data2": hasRoles,
	})
}
