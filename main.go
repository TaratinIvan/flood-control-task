package main

import (
	"context"
	"fmt"
	"github.com/TaratinIvan/flood-control-task/control"
	"github.com/go-redis/redis"
	"time"
)

// FloodControl interface definition moved to main.go
type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

func main() {
	// Redis connection details
	redisHost := "localhost"
	redisPort := 6379

	// Create a Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: "", // Set password if needed
		DB:       0,  // Use the default database
	})

	// Create a new flood control instance
	floodControl := control.NewFloodControl(redisClient, 10*time.Second, 5)

	// Check flood control for a user
	userID := int64(123)
	allowed, err := floodControl.Check(context.Background(), userID)
	if err != nil {
		fmt.Println("Error checking flood control:", err)
		return
	}

	// Print the result
	if allowed {
		fmt.Println("Request allowed for user:", userID)
	} else {
		fmt.Println("Request denied for user:", userID)
	}
}
