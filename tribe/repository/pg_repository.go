package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) tribe.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg tribe.FetchConfig) ([]*models.Tribe, int, error) {
	var err error
	data := []*models.Tribe{}
	total := 0
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Order(cfg.Sort...).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	if cfg.Filter != nil {
		query = query.Apply(cfg.Filter.Where)
	}

	if cfg.Count && cfg.Select {
		total, err = query.SelectAndCount()
	} else if cfg.Select {
		err = query.Select()
	} else if cfg.Count {
		total, err = query.Count()
	}
	if err != nil && err != pg.ErrNoRows {
		if strings.Contains(err.Error(), `relation "`+cfg.Server) {
			return nil, 0, fmt.Errorf("Server not found")
		}
		return nil, 0, fmt.Errorf("Internal server error")
	}

	return data, total, nil
}

func (repo *pgRepository) SearchTribe(ctx context.Context, cfg tribe.SearchTribeConfig) ([]*models.FoundTribe, int, error) {
	servers := []*models.Server{}
	if err := repo.
		Model(&servers).
		Context(ctx).
		Column("key").
		Where("version_code = ?", cfg.Version).
		Select(); err != nil {
		return nil, 0, fmt.Errorf("Internal server error")
	}

	var query *orm.Query
	res := []*models.FoundTribe{}
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
			Order(cfg.Sort...)
		if cfg.Count {
			count, err = base.SelectAndCount(&res)
		} else {
			err = base.Select(&res)
		}
		if err != nil && err != pg.ErrNoRows {
			return nil, 0, fmt.Errorf("Internal server error")
		}
	}

	return res, count, nil
}
