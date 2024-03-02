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

func (r RedisDB) GetOne(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	return &dto.Item{}, nil
}

func (r RedisDB) PutList(ctx context.Context, goods *dto.GetResponse) {
	jsoned, err := json.Marshal(goods)
	if err != nil {
		log.Warn().Err(err).Msg(fmt.Sprintf("redis: marshal %+v\n failed", goods.Meta))
	}
	key := strconv.Itoa(goods.Meta.Offset) + "-" + strconv.Itoa(goods.Meta.Limit)
	err = r.red.Set(ctx, key, jsoned, redisTTL*time.Second).Err()
	if err != nil {
		log.Warn().Err(err).Msg(fmt.Sprintf("redis: Put %+v\n failed", goods.Meta))
	}
}

func (r RedisDB) GetList(ctx context.Context, meta *dto.Meta) (string, error) {
	key := strconv.Itoa(meta.Offset) + "-" + strconv.Itoa(meta.Limit)

	val, err := r.red.Get(ctx, key).Result()
	if err != nil {
		log.Trace().Err(err).Msg(fmt.Sprintf("redis could not get list %s", key))
	}
	log.Trace().Msg(fmt.Sprintf("redis retrieved list %s", key))
	return val, err
}

func (r RedisDB) Invalidate(ctx context.Context, key string) {
	// Удаляем по ключю
	r.red.Del(ctx, key)
}
