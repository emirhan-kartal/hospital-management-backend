package handlers

import (
	"encoding/json"
	"example/hello/models"
	"example/hello/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PolyclinicHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (h *PolyclinicHandler) GetPolyclinicsNotInHospital(c *fiber.Ctx) error { //this is gonna be listed in a select / option tag in the frontend

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	hospitalID := claims["hospital_id"].(float64)

	var polyclinics []models.Polyclinic
	h.DB.Where("hospital_id = ?", hospitalID).Find(&polyclinics)

	var polyclinicsRedis []models.RedisPolyclinic

	polyclinicsString := h.Redis.Get(c.Context(), "polyclinics").Val()
	err := json.Unmarshal([]byte(polyclinicsString), &polyclinicsRedis)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println(polyclinicsRedis)
	fmt.Println(h.Redis.Get(c.Context(), "polyclinics").Val())
	var polyclinicsNotInHospital []string
	for i, polyclinic := range polyclinicsRedis {
		if len(polyclinics) == 0 {
			return c.JSON(polyclinicsRedis)
		}
		if polyclinics[i+1].Name != polyclinic.Name {
			polyclinicsNotInHospital = append(polyclinicsNotInHospital, polyclinic.Name)
		}
	}
	return c.JSON(polyclinicsNotInHospital)
}

type PolyclinicData struct {
	models.Polyclinic
	JobCounts []map[string]int
}

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

	fmt.Println(polyclinicData)

	for i := range polyclinics {
		jobTypeData := utils.BasicInfoHospital(&polyclinics[i])
		polyclinicData[i].JobCounts = append(polyclinicData[i].JobCounts, jobTypeData)
		polyclinicData[i].Polyclinic = polyclinics[i]
	}

	return c.JSON(polyclinicData)
}
func (h *PolyclinicHandler) AddPolyclinicToHospital(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}

	type Body struct {
		PolyclinicName string `json:"polyclinic_name"`
		City           string `json:"city"`
		District       string `json:"district"`
	}

	var body Body
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	hospitalID := claims["hospital_id"].(float64)

	var polyclinic models.Polyclinic

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !utils.RedisDataContains("polyclinics", body.PolyclinicName, c.Context(), h.Redis) { // not sure if this should be implemented
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Polyclinic",
		})
	}

	polyclinic.HospitalID = int(hospitalID)
	polyclinic.Name = body.PolyclinicName
	polyclinic.City = body.City
	polyclinic.District = body.District

	if err := h.DB.Create(&polyclinic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(polyclinic)
}

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
