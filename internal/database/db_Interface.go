package database

import (
	"context"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
)

type HezzlDbI interface {
	GoodCreateDB(ctx context.Context, item *dto.Item) (*dto.Item, error)
	GoodUpdateDB(ctx context.Context, item *dto.Item) (*dto.Item, error)
	GoodRemoveDB(ctx context.Context, info *dto.Item) (*dto.ItemShort, error)
	GoodsListDB(ctx context.Context, meta *dto.Meta) (*dto.GetResponse, error)
	GoodReprioritizeDB(ctx context.Context, item *dto.Item) (*dto.ReprResponse, error)
}
