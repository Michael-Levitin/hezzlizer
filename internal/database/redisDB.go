package database

import (
	"context"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/redis/go-redis/v9"
	"strconv"
)

const redisTTL = 60 // seconds

type RedisDB struct {
	red *redis.Client
}

func NewRedisDB(r *redis.Client) *RedisDB {
	return &RedisDB{red: r}
}

func (r RedisDB) Put(ctx context.Context, item *dto.Item) {

}

func (r RedisDB) Invalidate(ctx context.Context, item *dto.Item) {
	// Удаляем по ключю
	r.red.Del(ctx, strconv.Itoa(item.Id))
}

func (r RedisDB) Get(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	return &dto.Item{}, nil
}
