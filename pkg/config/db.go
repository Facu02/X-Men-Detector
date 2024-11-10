// config/redis.go
package config

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() (*redis.Client, error) {

	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("error al conectar con redis: %v", err)
		return redisClient, err
	}
	log.Println("conexion a redis exitosa")
	return redisClient, nil
}
