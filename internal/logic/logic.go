package logic

import (
	"context"
	"fmt"
	"github.com/Michael-Levitin/hezzlizer/internal/database"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/rs/zerolog/log"
)

type HezzlLogic struct {
	HezzlDB database.HezzlDbI
	RedisDB database.RedisDB
}

// NewHezzlLogic подключаем интерфейс БД в новую логику
func NewHezzlLogic(HezzlDB database.HezzlDbI, RedisDB *database.RedisDB) *HezzlLogic {
	return &HezzlLogic{HezzlDB: HezzlDB, RedisDB: *RedisDB}
}

func (h HezzlLogic) GoodCreate(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	if item.ProjectID == 0 || item.Name == "" {
		return &dto.Item{}, fmt.Errorf("projectId and name cannot be empty")
	}
	return h.HezzlDB.GoodCreateDB(ctx, item)
}

func (h HezzlLogic) GoodUpdate(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	if item.ProjectID == 0 || item.Name == "" || item.Id == 0 {
		return &dto.Item{}, fmt.Errorf("id, projectId & name cannot be empty")
	}
	h.RedisDB.Invalidate(ctx)
	return h.HezzlDB.GoodUpdateDB(ctx, item)
}

func (h HezzlLogic) GoodRemove(ctx context.Context, item *dto.Item) (*dto.ItemShort, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	if item.ProjectID == 0 || item.Id == 0 {
		return &dto.ItemShort{}, fmt.Errorf("id, projectId cannot be empty")
	}
	h.RedisDB.Invalidate(ctx)
	return h.HezzlDB.GoodRemoveDB(ctx, item)
}

func (h HezzlLogic) GoodsList(ctx context.Context, meta *dto.Meta) (*dto.GetResponse, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", meta))
	resp, err := h.HezzlDB.GoodsListDB(ctx, meta)
	if err != nil {
		return &dto.GetResponse{}, err
	}
	h.RedisDB.PutList(ctx, resp)
	return resp, err
}

func (h HezzlLogic) GoodReprioritize(ctx context.Context, item *dto.Item) (*dto.ReprResponse, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	if item.ProjectID == 0 || item.Id == 0 || item.Priority == 0 {
		return &dto.ReprResponse{}, fmt.Errorf("id, projectId & prioirity cannot be empty")
	}
	h.RedisDB.Invalidate(ctx)
	return h.HezzlDB.GoodReprioritizeDB(ctx, item)
}
