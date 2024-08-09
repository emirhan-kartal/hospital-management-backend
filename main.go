package main

import (
	"context"
	"example/hello/handlers"
	"example/hello/middleware"
	"example/hello/routes"
	"log"
	"os"

	_ "example/hello/docs"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var JWT_SECRET string
var ctx = context.Background()

// @title Auth API
// @version 1.0
// @description API for managing authentication
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db := handlers.SqliteHandler()
	redis := handlers.RedisHandler(ctx)

	gotoEnvErr := godotenv.Load()

	if gotoEnvErr != nil {
		log.Fatal("Error loading .env file")
	}
	JWT_SECRET = os.Getenv("JWT_SECRET")
	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	routes.SetupAuthRoutes(app, db, redis)
	routes.SetupPwResetRoutes(app, db, redis)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(JWT_SECRET)},
	}))

	app.Use(middleware.CheckTokenVersion(db))
	routes.SetupPersonelRoutes(app, db, redis)
	routes.SetupUserRoutes(app, db)
	routes.SetupPolyclinicRoutes(app, db, redis)

	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋! This is a protected route")
	})
	app.Listen(":3000")
}
