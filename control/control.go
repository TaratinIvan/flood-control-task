package control

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Интерфейс FloodControl
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

type redisFloodControl struct {
	redisClient *redis.Client
	maxRequests int
	windowSize  time.Duration
	mutex       sync.Mutex
}

func NewFloodControl(redisClient *redis.Client, maxRequests int, windowSize time.Duration) FloodControl {
	return &redisFloodControl{
		redisClient: redisClient,
		maxRequests: maxRequests,
		windowSize:  windowSize,
	}
}

func (f *redisFloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	key := fmt.Sprintf("flood_control:%d", userID)

	// Получение списка timestamp-ов
	timestamps, err := f.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return false, err
	}

	// Если за последние N секунд было совершено больше K вызовов
	if len(timestamps) >= f.maxRequests {
		// Вычисление разницы между текущим временем и timestamp-ом в начале списка
		oldestTimestamp, err := time.Parse(time.RFC3339, timestamps[0])
		if err != nil {
			return false, err
		}

		difference := time.Since(oldestTimestamp)

		// Если разница меньше N, то это означает, что лимит запросов превышен
		if difference < f.windowSize {
			return false, nil
		}
	}

	// Добавление timestamp-а в список
	err = f.redisClient.LPush(ctx, key, time.Now().Format(time.RFC3339)).Err()
	if err != nil {
		return false, err
	}

	// Ограничение длины списка
	f.redisClient.LTrim(ctx, key, 0, f.maxRequests-1)

	return true, nil
}
