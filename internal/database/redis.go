package database

import (
	"os"
	"log"
	"time"
	"math/rand"
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

	log.Println(redis_addr);

	maxRetryAttempts := 5;
	initialBackoff := 1 * time.Second

	for i := 0; i < maxRetryAttempts; i++ {

		rdb := redis.NewClient(&redis.Options{

			Addr:     redis_addr, 
			Password: "",
			DB:       0,
		})

		val, err := rdb.Ping(context.Background()).Result()
		if err == nil {

			log.Printf("Created Redis client...: %+v\n", val);
			return rdb
		}
			
		log.Printf("Error returned by Ping: %+v\nRetrying...\n", err)

		backoff := initialBackoff * (1 << i);
		jitter := time.Duration(rand.Int63n(int64(backoff) / 4));
		sleepDuration := backoff + jitter

		<-time.After(sleepDuration)
	}

	log.Println("Max retry attempts readched. Connecting to redis client failed.");

	return nil; 
}
