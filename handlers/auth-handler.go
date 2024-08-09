package handlers

import (
	"example/hello/models"
	"os"
	"time"

	_ "example/hello/docs"

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
type RegisterDataBody struct {
	Hospital models.HospitalBody
	User     models.UserBody
}

type LoginData struct {
	Email    string `json:"email"`
	Tel_no   string `json:"tel_no"`
	Password string `json:"password"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary Register a new user
// @Description Register a new user and hospital,Ignore HospitalID in User
// @Tags auth
// @Accept  json
// @Produce  json
// @Param registerData body RegisterDataBody true "Register Data"
// @Success 200 {object} RegisterData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
// @Security BearerAuth
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

	if data.User.Role != "Admin" && data.User.Role != "User" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role must be Admin or User",
		})
	}
	if err := h.DB.Create(&data.Hospital).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	data.User.HospitalID = data.Hospital.ID
	data.User.Password = string(hashedPassword)

	if err := h.DB.Create(&data.User).Error; err != nil {
		if err := h.DB.Unscoped().Delete(&data.Hospital).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	data.Hospital.User = append(data.Hospital.User, data.User)
	return c.JSON(data.Hospital)
}

// @Summary Login a user
// @Description Login a user and return a JWT token.You can login with email or tel_no
// @Tags auth
// @Accept  json
// @Produce  json
// @Param loginData body LoginData true "Login Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /login [post]
// @Security BearerAuth
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
