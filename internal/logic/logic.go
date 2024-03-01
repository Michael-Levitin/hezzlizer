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
}

// NewHezzlLogic подключаем интерфейс БД в новую логику
func NewHezzlLogic(HezzlDB database.HezzlDbI) *HezzlLogic {
	return &HezzlLogic{HezzlDB: HezzlDB}
}

func (h HezzlLogic) GoodCreate(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	if item.ProjectID == 0 || item.Name == "" {
		return &dto.Item{}, fmt.Errorf("projectId and name cannot be empty")
	}
	return h.HezzlDB.GoodCreate(ctx, item)
}

func (h HezzlLogic) GoodUpdate(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	if item.ProjectID == 0 || item.Name == "" || item.Id == 0 {
		return &dto.Item{}, fmt.Errorf("id, projectId & name cannot be empty")
	}
	return h.HezzlDB.GoodUpdate(ctx, item)
}

func (h HezzlLogic) GoodRemove(ctx context.Context, item *dto.Item) (*dto.ItemShort, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	if item.ProjectID == 0 || item.Id == 0 {
		return &dto.ItemShort{}, fmt.Errorf("id, projectId cannot be empty")
	}
	return h.HezzlDB.GoodRemove(ctx, item)
}

func (h HezzlLogic) GoodsList(ctx context.Context, info *dto.Meta) (*dto.GetResponse, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", info))
	return h.HezzlDB.GoodsList(ctx, info)
}

func (h HezzlLogic) GoodReprioritize(ctx context.Context, item *dto.Item) (*dto.ReprResponse, error) {
	log.Trace().Msg(fmt.Sprintf("Logic recieved %+v\n", item))
	return h.HezzlDB.GoodReprioritize(ctx, item)
}
