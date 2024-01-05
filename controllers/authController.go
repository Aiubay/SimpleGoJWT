package controllers

import (
	"fmt"
	"jwtreact/database"
	"jwtreact/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	// var data map[string]string

	var data models.User

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data.Passwords), 14)
	user := models.User{
		Name:      data.Name,
		Email:     data.Email,
		Passwords: password,
	}

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	// var data map[string]string
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data.Email).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Passwords, []byte(data.Passwords)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(), // 24 hours
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Cannot Login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Minute * 1),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": cookie,
	})

}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

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

	timeStart := time.Now().Unix()
	for i, row := range rows {
		for _, colCell := range row {
			exportsModel := models.Excel{
				Random: colCell,
			}
			database.DB.Create(&exportsModel)
		}

		fmt.Println(i + 1)
	}
	timeEnd := time.Now().Unix()

	return c.JSON(fiber.Map{
		"message":   "success",
		"totalTime": (timeEnd - timeStart) / 60,
	})

}
