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

	_goodRemoveQuery = `
UPDATE goods SET
       removed = true
WHERE id = @id AND project_id = @projectId
RETURNING id, project_id, removed;`

	_goodMetaQuery = `
SELECT
    (SELECT count(*) FROM goods) as total,
    (SELECT count(*) FROM goods WHERE removed = true) as removed;`

	_goodListQuery = `
SELECT id, project_id, name, COALESCE(description, '') as description, priority, removed, created_at
FROM goods
ORDER BY id
OFFSET @offset LIMIT @limit`

	_goodReprQuery = `
with updated as (
    UPDATE goods
        SET priority = priority + @priority - (
            SELECT priority
            FROM goods
            WHERE id = @id AND project_id = @projectId
        )

        WHERE priority >= (
            SELECT priority
            FROM goods
            WHERE id = @id AND project_id = @projectId
        )
        RETURNING id, priority)
select *
from updated
ORDER BY id;`
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

func (h HezzlDB) GoodRemove(ctx context.Context, item *dto.Item) (*dto.ItemShort, error) {
	log.Trace().Msg(fmt.Sprintf("DB recieve %+v\n", item))
	tx, err := h.db.Begin(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("failed beginning transaction"))
		return &dto.ItemShort{}, dto.ErrQueryExecute
	}
	itemS := dto.ItemShort{}
	err = h.db.QueryRow(ctx, _goodRemoveQuery,
		pgx.NamedArgs{"id": item.Id, "projectId": item.ProjectID}).
		Scan(&itemS.Id,
			&itemS.ProjectID,
			&itemS.Removed,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		tx.Commit(ctx)
		return &dto.ItemShort{}, fmt.Errorf("query found nothing to update")
	} else if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodRemove could not update %+v", item))
		err = tx.Rollback(ctx)
		if err != nil {
			log.Debug().Err(err).Msg(fmt.Sprintf("GoodRemove failed rolling back transaction"))
		}
		return &dto.ItemShort{}, dto.ErrQueryExecute
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodRemove failed commiting transaction"))
		return &dto.ItemShort{}, dto.ErrQueryExecute
	}

	return &itemS, nil
}

func (h HezzlDB) GoodsList(ctx context.Context, meta *dto.Meta) (*dto.GetResponse, error) {
	log.Trace().Msg(fmt.Sprintf("DB recieve %+v\n", meta))

	err := h.db.QueryRow(ctx, _goodMetaQuery).Scan(&meta.Total, &meta.Removed)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodsList could not set meta %+v", meta))
		return &dto.GetResponse{}, dto.ErrQueryExecute
	}

	rows, err := h.db.Query(ctx, _goodListQuery,
		pgx.NamedArgs{"limit": meta.Limit, "offset": meta.Offset})
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodsList could not get list %+v", meta))
		return &dto.GetResponse{}, dto.ErrQueryExecute
	}

	goods, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Item])
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("CollectRows error"))
		return &dto.GetResponse{}, dto.ErrQueryExecute
	}

	return &dto.GetResponse{
		Meta:  *meta,
		Goods: goods,
	}, nil
}

func (h HezzlDB) GoodReprioritize(ctx context.Context, item *dto.Item) (*dto.ReprResponse, error) {
	log.Trace().Msg(fmt.Sprintf("DB recieve %+v\n", item))

	tx, err := h.db.Begin(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("failed beginning transaction"))
		return &dto.ReprResponse{}, dto.ErrQueryExecute
	}

	rows, err := h.db.Query(ctx, _goodReprQuery,
		pgx.NamedArgs{"id": item.Id, "projectId": item.ProjectID, "priority": item.Priority})
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodReprioritize could not set priorities %+v", item))
		err = tx.Rollback(ctx)
		if err != nil {
			log.Debug().Err(err).Msg(fmt.Sprintf("GoodReprioritize failed rolling back transaction"))
		}
		return &dto.ReprResponse{}, dto.ErrQueryExecute
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("GoodReprioritize failed commiting transaction"))
		return &dto.ReprResponse{}, dto.ErrQueryExecute
	}

	priorities, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Priority])
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf("CollectRows error"))
		return &dto.ReprResponse{}, dto.ErrQueryExecute
	}

	return &dto.ReprResponse{Priorities: priorities}, nil
}
