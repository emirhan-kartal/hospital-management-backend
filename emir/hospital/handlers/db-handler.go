package handlers

import (
	"context"
	"emir/hospital/models"
	"emir/hospital/utils"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteHandler() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Hospital{}, &models.Personel{}, &models.Polyclinic{}, &models.User{}, &models.User{})
	fmt.Println("Connected to database")
	return db
}

func RedisHandler(ctx context.Context) *redis.Client {
	// redis handler
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}

	polyclinics := []models.RedisPolyclinic{
		{ID: 1, Name: "Kulak Burun Boğaz"},
		{ID: 2, Name: "Göz"},
		{ID: 3, Name: "Dahiliye"},
		{ID: 4, Name: "Kardiyoloji"},
		{ID: 5, Name: "Ortopedi"},
	}
	cities := []models.City{
		{ID: 1, Name: "İstanbul", Districts: []models.District{{ID: 1, Name: "Kadıköy"}, {ID: 2, Name: "Üsküdar"}}},
		{ID: 2, Name: "Ankara", Districts: []models.District{{ID: 1, Name: "Çankaya"}, {ID: 2, Name: "Keçiören"}}},
		{ID: 3, Name: "İzmir", Districts: []models.District{{ID: 1, Name: "Karşıyaka"}, {ID: 2, Name: "Konak"}}},
	}
	jobTitles := []models.JobType{
		{ID: 1, Name: "Doktor", Titles: []models.Title{{ID: 1, Name: "Asistan Doktor"}, {ID: 2, Name: "Uzman Doktor"}}},
		{ID: 2, Name: "Hemşire", Titles: []models.Title{{ID: 1, Name: "Uzman Hemşire"}, {ID: 2, Name: "Başhemşire"}}},
		{ID: 3, Name: "Teknisyen", Titles: []models.Title{{ID: 1, Name: "Laboratuvar Teknisyeni"}, {ID: 2, Name: "Radyoloji Teknisyeni"}}},
		{ID: 4, Name: "İdari Personel", Titles: []models.Title{{ID: 1, Name: "Başhekim"}, {ID: 2, Name: "Müdür"}, {ID: 3, Name: "Sekreter"}}},
		{ID: 5, Name: "Hizmet Personeli", Titles: []models.Title{{ID: 1, Name: "Temizlik Personeli"}, {ID: 2, Name: "Güvenlik Personeli"}}},
	}

	utils.SetData(ctx, rdb, "polyclinics", polyclinics)
	utils.SetData(ctx, rdb, "cities", cities)
	utils.SetData(ctx, rdb, "job_types", jobTitles)
	return rdb
}
