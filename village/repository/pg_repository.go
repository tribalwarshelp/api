package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) village.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg village.FetchConfig) ([]*models.Village, int, error) {
	var err error
	data := []*models.Village{}
	query := repo.WithParam("SERVER", pg.Safe(cfg.Server)).Model(&data).Context(ctx)

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter).
			Limit(cfg.Filter.Limit).
			Offset(cfg.Filter.Offset)

		if cfg.Filter.XGTE != 0 {
			query = query.Where("x >= ?", cfg.Filter.XGTE)
		} else if cfg.Filter.XGT != 0 {
			query = query.Where("x > ?", cfg.Filter.XGT)
		}
		if cfg.Filter.XLTE != 0 {
			query = query.Where("x <= ?", cfg.Filter.XLTE)
		} else if cfg.Filter.XLT != 0 {
			query = query.Where("x < ?", cfg.Filter.XLT)
		}

		if cfg.Filter.YGTE != 0 {
			query = query.Where("y >= ?", cfg.Filter.YGTE)
		} else if cfg.Filter.YGT != 0 {
			query = query.Where("y > ?", cfg.Filter.YGT)
		}
		if cfg.Filter.YLTE != 0 {
			query = query.Where("y <= ?", cfg.Filter.YLTE)
		} else if cfg.Filter.YLT != 0 {
			query = query.Where("y < ?", cfg.Filter.YLT)
		}

		if cfg.Filter.Sort != "" {
			query = query.Order(cfg.Filter.Sort)
		}

		if cfg.Filter.PlayerFilter != nil {
			query = query.Relation("Player._").WhereStruct(cfg.Filter.PlayerFilter)
			if cfg.Filter.PlayerFilter.TribeFilter != nil {
				query = query.
					Join("LEFT JOIN ?SERVER.tribes AS tribe ON tribe.id = player.tribe_id").
					WhereStruct(cfg.Filter.PlayerFilter.TribeFilter)
			}
		}
	}

	total := 0
	if cfg.Count {
		total, err = query.SelectAndCount()
	} else {
		err = query.Select()
	}
	if err != nil && err != pg.ErrNoRows {
		if strings.Contains(err.Error(), `relation "`+cfg.Server) {
			return nil, 0, fmt.Errorf("Server not found")
		}
		return nil, 0, errors.Wrap(err, "Internal server error")
	}

	return data, total, nil
}
