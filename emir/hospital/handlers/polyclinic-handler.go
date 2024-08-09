package handlers

import (
	"emir/hospital/models"
	"emir/hospital/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PolyclinicHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

type PolyclinicData struct {
	models.Polyclinic
	JobCounts []map[string]int
}

// @Summary Get polyclinics of the hospital
// @Description Retrieve a list of polyclinics associated with the user's hospital
// @Tags polyclinic
// @Accept json
// @Produce json
// @Success 200 {array} PolyclinicData
// @Failure 500 {object} ErrorResponse
// @Router /polyclinics [get]
// @Security BearerAuth
func (h *PolyclinicHandler) GetPolyclinicsOfHospital(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	hospitalID := claims["hospital_id"].(float64)
	var hospital models.Hospital
	if err := h.DB.Where("id = ?", hospitalID).Preload("Polyclinics").First(&hospital).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var polyclinics []models.Polyclinic

	if err := h.DB.Where("hospital_id = ?", hospitalID).Find(&polyclinics).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	polyclinicData := make([]PolyclinicData, len(polyclinics))

	for i := range polyclinics {
		jobTypeData := utils.BasicInfoHospital(&polyclinics[i])
		polyclinicData[i].JobCounts = append(polyclinicData[i].JobCounts, jobTypeData)
		polyclinicData[i].Polyclinic = polyclinics[i]
	}

	return c.JSON(polyclinicData)
}

// @Summary Add a polyclinic to the hospital
// @Description Add a new polyclinic to the user's hospital
// @Tags polyclinic
// @Accept json
// @Produce json
// @Param body body models.PolyclinicBody true "Polyclinic Data"
// @Success 200 {object} models.Polyclinic
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /polyclinics [post]
// @Security BearerAuth
func (h *PolyclinicHandler) AddPolyclinicToHospital(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}

	var body models.PolyclinicBody
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	hospitalID := claims["hospital_id"].(float64)

	var polyclinic models.Polyclinic

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !utils.RedisDataContains("polyclinics", body.Name, c.Context(), h.Redis) { // not sure if this should be implemented
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Polyclinic",
		})
	}

	polyclinic.HospitalID = int(hospitalID)
	polyclinic.Name = body.Name
	polyclinic.City = body.City
	polyclinic.District = body.District

	if err := h.DB.Create(&polyclinic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(polyclinic)
}

// @Summary Delete a polyclinic
// @Description Delete an existing polyclinic by ID
// @Tags polyclinic
// @Accept json
// @Produce json
// @Param id path int true "Polyclinic ID"
// @Success 200 {object} models.Polyclinic
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /polyclinics/{id} [delete]
// @Security BearerAuth
func (h *PolyclinicHandler) DeletePolyclinic(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}

	id := c.Params("id")

	var polyclinic models.Polyclinic
	if err := h.DB.Where("id = ?", id).First(&polyclinic).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Polyclinic not found",
		})
	}

	if err := h.DB.Unscoped().Delete(&polyclinic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
