package handlers

import (
	"example/hello/models"
	"os"
	"time"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

type RegisterData struct {
	Hospital models.Hospital
	User     models.User
}
type LoginData struct {
	Email    string `json:"email"`
	Tel_no   string `json:"tel_no"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var data RegisterData

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	data.User.Password = string(hashedPassword)
	fmt.Println("User:" + string(hashedPassword))

	if err := h.DB.Create(&data.Hospital).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	data.User.HospitalID = data.Hospital.ID

	if err := h.DB.Create(&data.User).Error; err != nil {
		if err := h.DB.Unscoped().Delete(&data.Hospital).Error; err != nil {
			fmt.Println("Hospital Silinemedi")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	data.Hospital.User = append(data.Hospital.User, data.User)
	return c.JSON(data)
}
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var user LoginData
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var dbUser models.User
	if user.Tel_no != "" {
		if err := h.DB.Where("tel_no = ?", user.Tel_no).First(&dbUser).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
	} else {
		if err := h.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	claims := jwt.MapClaims{
		"Ad Soyad":    dbUser.Name + " " + dbUser.Surname,
		"email":       dbUser.Email,
		"role":        dbUser.Role,
		"hospital_id": dbUser.HospitalID,
		"version":     dbUser.TokenVersion,
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	}
	jwt_secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwt_secret))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t})
}
