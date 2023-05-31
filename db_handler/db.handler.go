package db_handler

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	ID      uint   `gorm:"primary_key"`
	Site    string `gorm:"unique_index;not null"`
	Counter int64  `gorm:"not null" ;gorm:"default:0"`
}

func InitPostgres() (*gorm.DB, error) {
	POSTGRESHost := os.Getenv("POSTGRES_HOST")
	POSTGRESPort := os.Getenv("POSTGRES_PORT")
	POSTGRESUser := os.Getenv("POSTGRES_USER")
	POSTGRESPassword := os.Getenv("POSTGRES_PASSWORD")
	POSTGRESName := os.Getenv("POSTGRES_NAME")

	ConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s POSTGRESname=%s sslmode=disable parseTime=True, QueryFields: true", POSTGRESHost, POSTGRESPort, POSTGRESUser, POSTGRESPassword, POSTGRESName)

	db, err := gorm.Open(postgres.Open(ConnStr))
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Record{})
	return db, nil
}
func InitRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	return redisClient
}
