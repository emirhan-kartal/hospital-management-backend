package routes

import (
	"example/hello/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupPwResetRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {
	// Define your routes here
	app.Post("/pw-reset-initialize", (&handlers.PWHandler{DB: db, Redis: redis}).ResetPasswordInitiate)
	app.Post("/pw-reset-finalize", (&handlers.PWHandler{DB: db, Redis: redis}).ResetPasswordFinalize)
	app.Post("/pw-reset", (&handlers.PWHandler{DB: db, Redis: redis}).ResetPassword)

}
