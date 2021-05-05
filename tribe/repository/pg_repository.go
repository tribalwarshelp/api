package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/tribalwarshelp/api/tribe"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) tribe.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg tribe.FetchConfig) ([]*twmodel.Tribe, int, error) {
	var err error
	var data []*twmodel.Tribe
	total := 0
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset).
		Apply(cfg.Filter.Where).
		Apply(gopgutil.OrderAppender{
			Orders:   cfg.Sort,
			MaxDepth: 4,
		}.Apply)

	if cfg.Count && cfg.Select {
		total, err = query.SelectAndCount()
	} else if cfg.Select {
		err = query.Select()
	} else if cfg.Count {
		total, err = query.Count()
	}
	if err != nil && err != pg.ErrNoRows {
		if strings.Contains(err.Error(), `relation "`+cfg.Server) {
			return nil, 0, errors.New("Server not found")
		}
		return nil, 0, errors.New("Internal server error")
	}

	return data, total, nil
}

func (repo *pgRepository) SearchTribe(ctx context.Context, cfg tribe.SearchTribeConfig) ([]*twmodel.FoundTribe, int, error) {
	var servers []*twmodel.Server
	if err := repo.
		Model(&servers).
		Context(ctx).
		Column("key").
		Where("version_code = ?", cfg.Version).
		Select(); err != nil {
		return nil, 0, errors.New("Internal server error")
	}

	var query *orm.Query
	var res []*twmodel.FoundTribe
	for _, server := range servers {
		safeKey := pg.Safe(server.Key)
		otherQuery := repo.
			Model().
			Context(ctx).
			ColumnExpr("? AS server", server.Key).
			Column("tribe.id", "tribe.name", "tribe.tag", "tribe.most_points", "tribe.best_rank", "tribe.most_villages").
			TableExpr("?0.tribes as tribe", safeKey).
			Where("tribe.tag ILIKE ?0 OR tribe.name ILIKE ?0", cfg.Query)
		if query == nil {
			query = otherQuery
		} else {
			query = query.UnionAll(otherQuery)
		}
	}

	var err error
	count := 0
	if query != nil {
		base := repo.
			Model().
			With("union_q", query).
			Table("union_q").
			Limit(cfg.Limit).
			Offset(cfg.Offset).
			Apply(gopgutil.OrderAppender{
				Orders:   cfg.Sort,
				MaxDepth: 4,
			}.Apply)
		if cfg.Count {
			count, err = base.SelectAndCount(&res)
		} else {
			err = base.Select(&res)
		}
		if err != nil && err != pg.ErrNoRows {
			return nil, 0, errors.New("Internal server error")
		}
	}

	return res, count, nil
}
