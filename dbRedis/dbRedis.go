package dbRedis

import (
	"github.com/go-redis/redis"
	// "context"
	// redis "cloud.google.com/go/redis/apiv1"
)

func ConnectRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client

}

// func ConnectRedis() *redis.CloudRedisClient {

// 	ctx := context.Background()

// 	c, _ := redis.NewCloudRedisClient(ctx)

// 	return c

// }
