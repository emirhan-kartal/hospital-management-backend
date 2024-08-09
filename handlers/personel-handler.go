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
	LastId  int    `json:"last_id" query:"last_id" `
	Name    string `json:"name" query:"name"`
	Surname string `json:"surname" query:"surname"`
	TcNo    string `json:"tc_no" query:"tc_no"`
	JobType string `json:"job_type" query:"job_type"`
	Title   string `json:"title" query:"title"`
}
type Personel struct {
	TcNo        string `json:"tc_no"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	TelNo       string `json:"tel_no"`
	JobType     string `json:"job_type"`
	Title       string `json:"title"`
	WorkingDays string `json:"working_days"`
}

// @Summary Get personel list based on filters
// @Description Retrieve personel list based on filters.Made using cursor based pagination for better performance.If you put 0 to Last_ID it it will return the first 10 personels.
// @Tags personel
// @Accept json
// @Produce json
// @Param last_id query int false "Last ID"
// @Param name query string false "Name"
// @Param surname query string false "Surname"
// @Param tc_no query string false "TcNo"
// @Param job_type query string false "JobType"
// @Param title query string false "Title"
// @Success 200 {array} models.Personel
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /personels [get]
// @Security BearerAuth

func (h *PersonelHandler) GetPersonels(c *fiber.Ctx) error {

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)

	hospitalID := claims["hospital_id"].(float64)

	type Result struct {
		Personel
		PolyclinicName string `json:"polyclinic_name"`
	}

	var filters GetPersonelFilteringJson
	isFiltered := true

	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var results []Result
	query := h.DB.Table("personels").
		Select("personels.*, COALESCE(polyclinics.name, 'Belirsiz') as polyclinic_name").
		Joins("LEFT JOIN polyclinics ON polyclinics.id = personels.polyclinic_id").
		Where("personels.hospital_id = ? AND personels.id > ?", hospitalID, filters.LastId).
		Limit(10)

	if err := query.Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if filters.Name != "" {
		query = query.Where("personels.name LIKE ?", "%"+filters.Name+"%")
		isFiltered = true
	}
	if filters.Surname != "" {
		query = query.Where("personels.surname LIKE ?", "%"+filters.Surname+"%")
		isFiltered = true
	}
	if filters.TcNo != "" {
		query = query.Where("personels.tc_no LIKE ?", "%"+filters.TcNo+"%")
		isFiltered = true
	}
	if filters.JobType != "" {
		query = query.Where("personels.job_type LIKE ?", "%"+filters.JobType+"%")
		isFiltered = true
	}
	if filters.Title != "" {
		query = query.Where("personels.title LIKE ?", "%"+filters.Title+"%")
		isFiltered = true
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

// @Summary Get personel by ID
// @Description Retrieve a single personel by ID
// @Tags personel
// @Accept json
// @Produce json
// @Param id path int true "Personel ID"
// @Success 200 {object} models.Personel
// @Failure 404 {object} ErrorResponse
// @Router /personel/{id} [get]
// @Security BearerAuth
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

// @Summary Add a new personel
// @Description Add a new personel to the hospital
// @Tags personel
// @Accept json
// @Produce json
// @Param personel body models.PersonelBody true "personel"
// @Success 200 {object} models.Personel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /personel [post]
// @Security BearerAuth
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

type UpdatePersonelBody struct {
	Personel
	PolyclinicID int `json:"polyclinic_id"`
}

// @Summary Update a personel
// @Description Update an existing personel's details.gorm.Model shouldnt be included here.(createdAt,updatedAt,deletedAt)
// @Tags personel
// @Accept json
// @Produce json
// @Param id path int true "personel ID"
// @Param personel body UpdatePersonelBody true "personel"
// @Success 200 {object} models.Personel
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /personel/{id} [put]
// @Security BearerAuth
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
	if err := c.BodyParser(&personelBody); err != nil {
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

	if err := h.DB.Model(&personel).Updates(personelBody).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(personel)
}

// @Summary Delete a personel
// @Description Delete an existing personel by ID
// @Tags personel
// @Accept json
// @Produce json
// @Param id path int true "Personel ID"
// @Security BearerAuth
// @Success 200 {object} models.Personel
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /personel/{id} [delete]
func (h *PersonelHandler) DeletePersonel(c *fiber.Ctx) error {
	if !utils.IsAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not enough permission",
		})
	}
	id := c.Params("id")
	var personel models.Personel
	if err := h.DB.Table("personels").Where("id = ?", id).First(&personel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Personel not found",
		})
	}

	if err := h.DB.Unscoped().Delete(&personel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(personel)
}

func validatePersonel(p models.Personel, DB gorm.DB, Redis *redis.Client, c *fiber.Ctx) bool {
	if p.Title == "Başhekim" {
		var count int64
		DB.Model(&models.Personel{}).Where("title = ?", "Başhekim").Count(&count)
		if count >= 1 {
			return false
		}
	}
	redisDataJobTypes := Redis.Get(c.Context(), "job_types").Val()
	if !utils.TitleJobTypeRelated(p.Title, p.JobType, Redis) {
		return false
	}

	if !strings.Contains(redisDataJobTypes, p.Title) || !strings.Contains(redisDataJobTypes, p.JobType) {
		fmt.Println("Invalid Title or Job Type")
		return false
	}

	return true

}
