package handlers

import (
	"emir/hospital/models"
	"emir/hospital/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

// @Summary Create a new user
// @Description Add a new user to the hospital.HospitalID is taken from the token.It will be ignored in the request body
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UserBody true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
// @Security BearerAuth
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
	if user.Role != "Admin" && user.Role != "User" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Role",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.HospitalID = uint(HospitalID)

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

// @Summary Get all users
// @Description Retrieve a list of all users associated with the user's hospital
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
// @Security BearerAuth
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

// @Summary Get user by ID
// @Description Retrieve a single user by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
// @Security BearerAuth
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

// @Summary Update a user
// @Description Update an existing user's details
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UserBody true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
// @Security BearerAuth
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
	if err := h.DB.Unscoped().Delete(&targetUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
