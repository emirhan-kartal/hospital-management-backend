package middleware

import (
	"emir/hospital/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func CheckTokenVersion(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		tokenVersion := int(claims["version"].(float64))

		var dbUser models.User
		if err := db.Where("email = ?", email).First(&dbUser).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}

		if dbUser.TokenVersion != tokenVersion {
			return c.Status(fiber.StatusUnauthorized).SendString("Token version mismatch")
		}

		return c.Next()
	}
}
