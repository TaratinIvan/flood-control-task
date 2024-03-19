package control

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// FloodControl implementation (no interface definition here)
type floodControl struct {
	redisClient *redis.Client
	mutex       sync.Mutex
	interval    time.Duration
	limit       int
}

// NewFloodControl constructor (same as before)
func NewFloodControl(redisClient *redis.Client, interval time.Duration, limit int) *floodControl {
	return &floodControl{
		redisClient: redisClient,
		mutex:       sync.Mutex{},
		interval:    interval,
		limit:       limit,
	}
}

// Check method implementation (same as before)
func (f *floodControl) Check(ctx context.Context, userID int64) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	key := strconv.FormatInt(userID, 10)

	// Get timestamp and count from Redis
	timestamp, err := f.redisClient.Get(key).Int64()
	if err != nil {
		return false, err
	}

	count, err := f.redisClient.Incr(key).Result()
	if err != nil {
		return false, err
	}

	// Check time since last request
	now := time.Now()
	elapsed := now.Sub(time.Unix(timestamp, 0))

	// Reset count if interval has passed
	if elapsed >= f.interval {
		f.redisClient.Set(key, now.Unix(), 0)
		count = 1
	}

	// Return false if count exceeds limit
	return count <= int64(f.limit), nil
}
