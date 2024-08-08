package handlers

import (
	"example/hello/models"
	"example/hello/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}
	var senderUser = c.Locals("user").(*jwt.Token)
	claims := senderUser.Claims.(jwt.MapClaims)
	HospitalID := claims["hospital_id"].(float64)

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.HospitalID = int(HospitalID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user.Password = string(hashedPassword)
	if err := h.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)

}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	HospitalID := claims["hospital_id"].(float64)
	var users []models.User
	if err := h.DB.Where("hospital_id = ?", HospitalID).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(users)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var targetUser models.User
	if err := h.DB.Where("id = ?", id).First(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(targetUser)
}
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}
	id := c.Params("id")
	var targetUser models.User
	if err := h.DB.Where("id = ?", id).First(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if targetUser.ID == 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cannot update super admin",
		})
	}

	if err := c.BodyParser(&targetUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := h.DB.Save(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(targetUser)
}
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}
	id := c.Params("id")
	var targetUser models.User
	if err := h.DB.Where("id = ?", id).First(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err := h.DB.Delete(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
