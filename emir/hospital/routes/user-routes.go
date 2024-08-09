package routes

import (
	"emir/hospital/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB) {
	// Define your routes here
	// app.Post("/register", (&handlers.AuthHandler{DB: db}).Register)
	// app.Post("/login", (&handlers.AuthHandler{DB: db}).Login)
	app.Get("/users", (&handlers.UserHandler{DB: db}).GetUsers)
	app.Get("/users/:id", (&handlers.UserHandler{DB: db}).GetUser)
	app.Post("/users", (&handlers.UserHandler{DB: db}).CreateUser)
	app.Put("/users", (&handlers.UserHandler{DB: db}).UpdateUser)
	app.Delete("/users", (&handlers.UserHandler{DB: db}).DeleteUser)

}
