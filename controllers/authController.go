package controllers

import (
	"errors"
	"fmt"
	"jwtreact/database"
	"jwtreact/models"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var payload models.PayloadRegister
	var role models.Role

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := database.DB.Where("id = ?", 2).First(&role).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Fetch roles error",
			"err":     err,
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: password,
	}

	if result := tx.Create(&user); result.Error != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": result.Error,
		})
	}

	tx.Commit()

	if !AssignRole(int(user.ID), "guest") {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Failed to assign role",
		})
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	input := new(models.PayloadLogin)
	if err := c.BodyParser(input); err != nil {
		return err
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return err
	}

	if errCompare := bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password)); errCompare != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Incorrect password",
			"error":   errCompare,
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"claims":  claims,
	})
}

func GetUser(c *fiber.Ctx) error {
	jwtCookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(jwtCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	if claims.ExpiresAt < time.Now().Unix() {
		return c.JSON(fiber.Map{
			"message": "Claims expired",
		})
	}
	var user models.User
	database.DB.Where("id = ?", claims.Issuer).Preload("Role").First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Upload(c *fiber.Ctx) error {
	f, err := excelize.OpenFile("Book1.xlsx")

	if err != nil {
		return err
	}

	rows, err := f.GetRows("Sheet1")

	if err != nil {
		fmt.Println(err)
		return err
	}

	exportsModels := make([]models.Excel, 0, len(rows)*len(rows[0]))

	for _, row := range rows {
		for _, colCell := range row {
			exportsModels = append(exportsModels, models.Excel{
				Random: colCell,
			})
		}
	}

	if err := database.DB.Create(&exportsModels).Error; err != nil {
		return err
	}

	timeStart := time.Now()

	return c.JSON(fiber.Map{
		"message":   "success",
		"totalTime": time.Since(timeStart).Seconds(),
	})
}

func Testing(c *fiber.Ctx) error {
	// passwordRequest := "Testx"
	passwordDB := "Testx"

	bsp, err := bcrypt.GenerateFromPassword([]byte(passwordDB), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(bsp)
	new := []byte("$2a$10$Fr/MwLaF.VDPKbFyOPIMmukKi9mog9BhmHr07ASXVLJxKZTfQLCoS")
	err = bcrypt.CompareHashAndPassword(new, []byte("LongAssPasswordss"))
	if err != nil {
		panic(err)
	} else {
		fmt.Println("password are equal")
	}

	return nil
}
