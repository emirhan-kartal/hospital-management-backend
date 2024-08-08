package routes

import (
	"example/hello/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {
	// Define your routes here
	app.Post("/register", (&handlers.AuthHandler{DB: db}).Register)
	app.Post("/login", (&handlers.AuthHandler{DB: db}).Login)
}
