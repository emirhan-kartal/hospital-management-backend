package routes

import (
	"example/hello/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupPersonelRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {
	app.Get("/personels", (&handlers.PersonelHandler{DB: db, Redis: redis}).GetPersonels)
	app.Post("/personel/add", (&handlers.PersonelHandler{DB: db, Redis: redis}).AddPersonel)
	app.Put("/personel/update/:id", (&handlers.PersonelHandler{DB: db, Redis: redis}).UpdatePersonel)

	/* 	app.Get("/personel/:id", (&handlers.PersonelHandler{DB: db, Redis: redis}).GetPersonelByID)
	   	app.Delete("/personel/delete", (&handlers.PersonelHandler{DB: db, Redis: redis}).DeletePersonel) */

}
