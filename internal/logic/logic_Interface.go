package logic

import (
	"context"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
)

type HezzlLogicI interface {
	GoodCreate(ctx context.Context, item *dto.Item) (*dto.Item, error)
	GoodUpdate(ctx context.Context, item *dto.Item) (*dto.Item, error)
	GoodRemove(ctx context.Context, info *dto.Item) (*dto.Item, error)
	GoodsList(ctx context.Context, item *dto.Meta) (*dto.GetResponse, error)
	GoodReprioritize(ctx context.Context, item *dto.Item) (*dto.ReprResponse, error)
}
