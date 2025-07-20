package database

import (
	"os"
	"log"
	"context"
	"github.com/redis/go-redis/v9"
)

const (
	DRIVER_HASH = "driverLocationHash"
)

func NewRedisClient() *redis.Client {

	redis_addr := os.Getenv("REDIS_ADDR");
	if redis_addr == "" {

		redis_addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{

		Addr:     redis_addr, 
		Password: "",
		DB:       0,
	})

	val, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		
		log.Printf("Error returned by Ping: %w\n", err)
	}

	log.Println("Created Redis client...", val);

	return rdb
}
