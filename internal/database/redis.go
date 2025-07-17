package database

import (
	"log"
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {


	rdb := redis.NewClient(&redis.Options{

		Addr:     "localhost:6379",
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
