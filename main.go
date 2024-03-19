package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"

	"TaratinIvan/flood-control-task/control"
)

func main() {
	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Создание экземпляра FloodControl
	floodControl := control.NewFloodControl(redisClient, 10, 5)

	// Контекст
	ctx := context.Background()

	// Имитация вызовов Check
	for i := 0; i < 20; i++ {
		userID := xid.New().String()

		allowed, err := floodControl.Check(ctx, userID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("UserID: %s, Allowed: %v\n", userID, allowed)

		time.Sleep(time.Second)
	}
}
