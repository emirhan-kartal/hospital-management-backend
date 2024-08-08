package utils

import (
	"context"
	"encoding/json"
	"example/hello/models"
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

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
	fmt.Println(role)
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
		fmt.Println(v.Name + "  " + value)
		if v.Name == value {
			fmt.Println("True " + v.Name + "  " + value)
			return true
		}
	}
	return false
}
func getJobTypeCounts(db *gorm.DB, hospitalID int) (map[string]int, error) {
	var result []struct {
		JobType string
		Count   int
	}

	err := db.Table("personels").
		Select("job_type, COUNT(*) as count").
		Joins("JOIN polyclinics ON personels.polyclinic_id = polyclinics.id").
		Where("polyclinics.hospital_id = ?", hospitalID).
		Group("job_type").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	jobTypeCounts := make(map[string]int)
	for _, row := range result {
		jobTypeCounts[row.JobType] = row.Count
	}

	return jobTypeCounts, nil
}
