package controllers

import (
	"jwtreact/database"
	"jwtreact/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

// GetUsers retrieves a list of users from the database and returns it as JSON.
//
// Parameter(s): c *fiber.Ctx
// Return type(s): error
func GetUsers(c *fiber.Ctx) error {
	var users []models.User

	result := database.DB.Preload(clause.Associations).Find(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
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
		"data":   users,
	})
}
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id") // assumes that id is passed as a parameter in the URL
	var user models.User

	result := database.DB.Preload(clause.Associations).First(&user, id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "OK",
		"data":   user,
	})
}
