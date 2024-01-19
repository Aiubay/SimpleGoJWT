package controllers

import (
	"fmt"
	"jwtreact/database"
	"jwtreact/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {

	var data models.PayloadRegister
	var role models.Role
	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	tx := database.DB.Begin()

	resultRoles := database.DB.Where("id = ?", 2).First(&role)

	if resultRoles.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Fetch roles err",
			"err":     resultRoles.Error,
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: password,
	}

	result := tx.Create(&user)

	if result.Error != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": result.Error,
		})
	}
	tx.Commit()

	assignRole := AssignRole(int(user.ID), "guest")

	if !assignRole {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Failed to assign role",
		})
	}
	tx.Commit()
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {

	data := new(models.PayloadLogin)

	if err := c.BodyParser(data); err != nil {
		return err
	}

	var user models.User
	database.DB.Where("email = ?", data.Email).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	errCompare := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password))
	if errCompare != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
			"error":   errCompare,
		})
	}
	// checkCookie := c.Cookies("jwt")

	// if checkCookie != "" {
	// 	claims, status := extractClaims(checkCookie)
	// 	if status {
	// 		x := claims.ExpiresAt

	// 		if x < time.Now().Unix() {
	// 			return c.JSON(fiber.Map{
	// 				"status": "Cookies Checked",
	// 			})
	// 		}

	// 	}

	// }

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
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

// func extractClaims(tokenStr string) (*jwt.StandardClaims, bool) {
// 	hmacSecretString := SecretKey
// 	hmacSecret := []byte(hmacSecretString)

// 	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		// Check the signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		// Return the secret key
// 		return hmacSecret, nil
// 	})

// 	if err != nil {
// 		return nil, false
// 	}

// 	// Extract the claims
// 	claims, ok := token.Claims.(*jwt.StandardClaims)
// 	if !ok {
// 		return nil, false
// 	}

// 	return claims, true
// }

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

	if claims.ExpiresAt < time.Now().Unix() {
		return c.JSON(fiber.Map{
			"message": "Claims expired",
		})
	}
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
