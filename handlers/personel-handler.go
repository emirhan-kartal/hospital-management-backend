package handlers

import (
	"example/hello/models"
	"example/hello/utils"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PersonelHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}
type GetPersonelFilteringJson struct {
	LastId  int    `json:"last_id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	TcNo    string `json:"tc_no"`
	JobType string `json:"job_type"`
	Title   string `json:"title"`
}

func (h *PersonelHandler) GetPersonels(c *fiber.Ctx) error {

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)

	hospitalID := claims["hospital_id"].(float64)

	type Personel struct {
		TcNo           string `json:"tc_no"`
		Name           string `json:"name"`
		Surname        string `json:"surname"`
		TelNo          string `json:"tel_no"`
		JobType        string `json:"job_type"`
		Title          string `json:"title"`
		WorkingDays    string `json:"working_days"`
		PolyclinicName string `json:"polyclinic_name"`
	}

	type Result struct {
		Personel
		PolyclinicName string `json:"polyclinic_name"`
	}

	var filters GetPersonelFilteringJson
	isFiltered := true
	if err := c.BodyParser(&filters); err != nil {
		filters.LastId = 0
		isFiltered = false
	}
	var results []Result
	query := h.DB.Table("personels").
		Select("personels.*, polyclinics.name as polyclinic_name").
		Joins("JOIN polyclinics ON polyclinics.id = personels.polyclinic_id").
		Where("personels.hospital_id = ? AND personels.id > ?", hospitalID, filters.LastId).
		Limit(10)

	if err := query.Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if filters.Name != "" {
		query = query.Where("personels.name LIKE ?", "%"+filters.Name+"%")
	}
	if filters.Surname != "" {
		query = query.Where("personels.surname LIKE ?", "%"+filters.Surname+"%")
	}
	if filters.TcNo != "" {
		query = query.Where("personels.tc_no LIKE ?", "%"+filters.TcNo+"%")
	}
	if filters.JobType != "" {
		query = query.Where("personels.job_type LIKE ?", "%"+filters.JobType+"%")
	}
	if filters.Title != "" {
		query = query.Where("personels.title LIKE ?", "%"+filters.Title+"%")
	}

	if isFiltered {
		err := query.Find(&results).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.JSON(results)
}
func (h *PersonelHandler) GetPersonelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var personel models.Personel
	if err := h.DB.Where("id = ?", id).First(&personel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Personel not found",
		})
	}
	return c.JSON(personel)
}
func (h *PersonelHandler) AddPersonel(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}
	hospitalID := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["hospital_id"].(float64)

	var personel models.Personel
	if err := c.BodyParser(&personel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	personel.HospitalID = int(hospitalID)

	if !validatePersonel(personel, *h.DB, h.Redis, c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Title or Job Type or Baş Hekim already exists",
		})
	}
	if err := h.DB.Create(&personel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(personel)
}

func (h *PersonelHandler) UpdatePersonel(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}
	id := c.Params("id")
	var personel models.Personel
	if err := h.DB.Where("id = ?", id).First(&personel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Personel not found",
		})
	}
	var personelBody models.Personel
	if err := c.BodyParser(&personel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if personelBody.ID != 0 || personelBody.HospitalID != 0 {
		if (personel.ID != personelBody.ID) || (personel.PolyclinicID != personelBody.PolyclinicID || personel.HospitalID != personelBody.HospitalID) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Personel ID or Polyclinic ID cannot be changed",
			})
		}
	}

	if !validatePersonel(personelBody, *h.DB, h.Redis, c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Title or Job Type or Baş Hekim already exists",
		})
	}
	if err := h.DB.Save(&personel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(personel)
}

func validatePersonel(p models.Personel, DB gorm.DB, Redis *redis.Client, c *fiber.Ctx) bool {
	if p.Title == "Baş Hekim" {
		var count int64
		DB.Model(&models.Personel{}).Where("title = ?", "Baş Hekim").Count(&count)
		if count >= 1 {

			fmt.Println("Baş Hekim already exists")
			return false
		}
	}
	redisDataJobTypes := Redis.Get(c.Context(), "job_types").Val()
	fmt.Println((redisDataJobTypes))
	if !strings.Contains(redisDataJobTypes, p.Title) || !strings.Contains(redisDataJobTypes, p.JobType) {
		fmt.Println("Invalid Title or Job Type")
		return false
	}

	return true

}
