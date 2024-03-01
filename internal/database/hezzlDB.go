package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Michael-Levitin/hezzlizer/internal/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	_goodCreateQuery = `
INSERT INTO goods (project_id, name, description, priority, removed, created_at)
VALUES (@projectId, @name, NULL, (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods), false, NOW())
RETURNING id, project_id, name, COALESCE(description,''), priority, removed, created_at;`

	_goodUpdateQuery = `
UPDATE goods SET
       name = @name,
       description = COALESCE(@description, description)
WHERE id = @id AND project_id = @projectId
RETURNING id, project_id, name, COALESCE(description,''), priority, removed, created_at;`
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
		return &dto.Item{}, dto.ErrQueryExecute
	}

	return item, nil
}

func (h HezzlDB) GoodUpdate(ctx context.Context, item *dto.Item) (*dto.Item, error) {
	log.Trace().Msg(fmt.Sprintf("DB recieve %+v\n", item))
	tx, err := h.db.Begin(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("failed beginning transaction"))
		return &dto.Item{}, dto.ErrQueryExecute
	}

	err = h.db.QueryRow(ctx, _goodUpdateQuery,
		pgx.NamedArgs{"id": item.Id, "projectId": item.ProjectID, "name": item.Name, "description": item.Description}).
		Scan(&item.Id,
			&item.ProjectID,
			&item.Name,
			&item.Description,
			&item.Priority,
			&item.Removed,
			&item.CreatedAt,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		tx.Commit(ctx)
		return &dto.Item{}, fmt.Errorf("query found nothing to update")
	} else if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodUpdate could not update %+v", item))
		err = tx.Rollback(ctx)
		if err != nil {
			log.Debug().Err(err).Msg(fmt.Sprintf("GoodUpdate failed rolling back transaction"))
		}
		return &dto.Item{}, dto.ErrQueryExecute
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodUpdate failed commiting transaction"))
		return &dto.Item{}, dto.ErrQueryExecute
	}

	return item, nil
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

