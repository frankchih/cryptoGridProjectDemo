package redisLib

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisService struct {
	Rdb *redis.Client
}

var (
	TASK_QUOTE = "TaskQuote"
	TASK_ORDER = "TaskOrder"
)

func NewRedisService(rdb *redis.Client) *RedisService {
	return &RedisService{Rdb: rdb}
}

func (redisService *RedisService) SetSymbolPrice(symbol, price string) error {
	rdb := redisService.Rdb
	key := "PRICE_" + symbol
	value := price
	err := rdb.Set(context.Background(), key, value, 5*time.Second).Err()
	return err
}

func (redisService *RedisService) GetTaskHearthBeat(task string) (string, error) {
	rdb := redisService.Rdb
	key := "HEARTH_BEAT_" + task
	val, err := rdb.Get(context.Background(), key).Result()
	return val, err
}
func (redisService *RedisService) GetTTLTaskHearthBeat(task string) (time.Duration, error) {
	rdb := redisService.Rdb
	key := "HEARTH_BEAT_" + task
	val, err := rdb.TTL(context.Background(), key).Result()
	return val, err
}
func (redisService *RedisService) SetTaskHearthBeat(task string) error {
	rdb := redisService.Rdb
	key := "HEARTH_BEAT_" + task
	err := rdb.Set(context.Background(), key, key, 15*time.Second).Err()
	return err
}
func (redisService *RedisService) DelTaskHearthBeat(task string) error {
	rdb := redisService.Rdb
	key := "HEARTH_BEAT_" + task
	err := rdb.Del(context.Background(), key).Err()

	return err
}
