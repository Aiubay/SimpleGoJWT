package routes

import (
	"jwtreact/controllers"
	"jwtreact/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api", middleware.JwtMiddleware)

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	api.Post("/logout", controllers.Logout)

	api.Get("/user", controllers.GetUser)
	api.Get("/users", controllers.GetUsers)
	api.Get("/userById/:id", controllers.GetUserByID)

	api.Get("/getAllRole", controllers.GetAllRole)
	api.Get("/insertRole", controllers.InsertRole)
	api.Get("/deleteUsers", controllers.DeleteAllUsers)
	api.Get("/getUserById", controllers.GetUserByID)
	api.Get("/createUsers", controllers.CreateUsers)

	api.Get("/upload", controllers.Upload)

	app.Get("/testing", controllers.Testing)
}
