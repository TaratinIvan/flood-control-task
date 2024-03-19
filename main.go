package main

import (
	"context"
	"fmt"
	"time"

	"github.com/TaratinIvan/flood-control-task/control"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
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

// Интерфейс FloodControl
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
