package handlers

import (
	"emir/hospital/models"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RedisData struct {
	Redis *redis.Client
}

func (h *RedisData) GetJobTypes(c *fiber.Ctx) error {
	jobTypesString := h.Redis.Get(c.Context(), "job-types").Val()
	var jobTypes []models.JobType
	err := json.Unmarshal([]byte(jobTypesString), &jobTypes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(jobTypes)
}
