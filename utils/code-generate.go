package utils

import (
	"context"
	"encoding/json"
	"example/hello/models"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/golang-jwt/jwt/v5"
)

func RandCode(n int) string {
	var letters = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func RandString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ#!$%&/()=?")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func IsAdmin(c *fiber.Ctx) bool {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "Admin" {
		return false
	}
	return true

}

func BasicInfoHospital(h *models.Polyclinic) map[string]int {
	personelJobTypeCountInfo := make(map[string]int)
	for _, personel := range h.Personels {
		personelJobTypeCountInfo[personel.JobType]++
	}
	return personelJobTypeCountInfo

}

func SetData(ctx context.Context, rdb *redis.Client, key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
func RedisDataContains(key string, value string, ctx context.Context, rdb *redis.Client) bool {
	jsonData := rdb.Get(ctx, key).Val()
	var slice []models.RedisPolyclinic
	err := json.Unmarshal([]byte(jsonData), &slice)
	if err != nil {
		return false
	}
	for _, v := range slice {
		if v.Name == value {
			return true
		}
	}
	return false
}
func TitleJobTypeRelated(title string, jobType string, rdb *redis.Client) bool {
	var slice []models.JobType
	err := json.Unmarshal([]byte(rdb.Get(context.Background(), "job-types").Val()), &slice)
	if err != nil {
		return false
	}
	for _, v := range slice {
		for _, v_title := range v.Titles {
			if v_title.Name == title && v.Name == jobType {
				return true
			}
		}
	}

	return false
}
