package dbRedis

import (
	"os"
	"github.com/go-redis/redis"
)

func ConnectRedis() *redis.Client {
	REDIS_IP_PORT := os.Getenv("REDIS_IP_PORT")
	client := redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		//Password: "",
		Addr:	REDIS_IP_PORT,
		DB:   0,
	})

	return client

}
