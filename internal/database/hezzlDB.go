package database

import (
	"context"
	"fmt"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	_goodCreateQuery = `
INSERT INTO goods (project_id, name, description, priority, removed, created_at)
VALUES (@projectId, @name, '', (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods), false, NOW())
RETURNING id, project_id, name, description, priority, removed, created_at;`
)

type HezzlDB struct {
	db *pgxpool.Pool
}

func NewHezzlDB(db *pgxpool.Pool) *HezzlDB {
	return &HezzlDB{db: db}
}

func (h HezzlDB) GoodCreate(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	log.Trace().Msg(fmt.Sprintf("DB recieve %+v\n", item))

	err := h.db.QueryRow(ctx, _goodCreateQuery,
		pgx.NamedArgs{"projectId": item.ProjectID, "name": item.Name, "description": item.Description}).
		Scan(&item.Id,
			&item.ProjectID,
			&item.Name,
			&item.Description,
			&item.Priority,
			&item.Removed,
			&item.CreatedAt,
		)
	if err != nil {
		fmt.Println(err)
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodCreate could not create %+v", item))
		return &dto.Item{}, dto.QueryExecuteErorr
	}

	return item, nil
}

func (h HezzlDB) GoodUpdate(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (h HezzlDB) GoodRemove(ctx context.Context, info *dto.Item) (*dto.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (h HezzlDB) GoodsList(ctx context.Context, item *dto.Meta) (*dto.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h HezzlDB) GoodReprioritize(ctx context.Context, item *dto.Item) (*dto.ReprResponse, error) {
	//TODO implement me
	panic("implement me")
}

