package routes

import (
	"emir/hospital/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupPwResetRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {
	// Define your routes here
	app.Post("/reset-password/initiate", (&handlers.PWHandler{DB: db, Redis: redis}).ResetPasswordInitiate)
	app.Post("/reset-password/finalize", (&handlers.PWHandler{DB: db, Redis: redis}).ResetPasswordFinalize)
	app.Post("/reset-password", (&handlers.PWHandler{DB: db, Redis: redis}).ResetPassword)

}
