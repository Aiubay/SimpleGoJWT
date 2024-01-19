package controllers

import (
	"jwtreact/database"
	"jwtreact/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllRole(c *fiber.Ctx) error {

	var roles []models.Role

	result := database.DB.Find(&roles)

	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Fetch roles error",
			"Error":   result.Error,
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"data":    "No Data",
			"message": "Success",
		})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    roles,
		"message": "Success",
	})
}


func AssignRole(userId int, role string) bool {

	var user models.User

	findUser := database.DB.Where("id=?", userId).First(&user)

	if findUser.Error != nil || findUser.RowsAffected == 0 {
		return false
	}
	var roles models.Role

	findRole := database.DB.Where("slug=?", role).First(&roles)

	if findRole.Error != nil || findRole.RowsAffected == 0 {
		return false
	}

	user.Role = models.UserHasRoles{
		Role: roles.ID,
	}

	database.DB.Updates(&user)

	return true
}
