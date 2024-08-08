package routes

import (
	"example/hello/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupPolyclinicRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {
	app.Get("/polyclinics", (&handlers.PolyclinicHandler{DB: db, Redis: redis}).GetPolyclinicsOfHospital)
	app.Get("/polyclinics/not-in-hospital", (&handlers.PolyclinicHandler{DB: db, Redis: redis}).GetPolyclinicsNotInHospital)
	app.Post("/polyclinics/add", (&handlers.PolyclinicHandler{DB: db, Redis: redis}).AddPolyclinicToHospital)
	app.Post("/polyclinics/delete", (&handlers.PolyclinicHandler{DB: db, Redis: redis}).DeletePolyclinic)
}
