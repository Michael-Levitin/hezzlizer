package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

const redisTTL = 60 // seconds

type RedisDB struct {
	red *redis.Client
}

func NewRedisDB(r *redis.Client) *RedisDB {
	return &RedisDB{red: r}
}

func (r RedisDB) PutOne(ctx context.Context, item *dto.Item) {
	jsoned, err := json.Marshal(item)
	if err != nil {
		log.Warn().Err(err).Msg(fmt.Sprintf("redis: marshal %+v\n failed", item))
	}

	err = r.red.Set(ctx, strconv.Itoa(item.Id), jsoned, redisTTL*time.Second).Err()
	if err != nil {
		log.Warn().Err(err).Msg(fmt.Sprintf("redis: Put %+v\n failed", item))
	}
}

func (r RedisDB) Invalidate(ctx context.Context, key string) {
	// Удаляем по ключю
	r.red.Del(ctx, key)
}

func (r RedisDB) Get(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	return &dto.Item{}, nil
}
