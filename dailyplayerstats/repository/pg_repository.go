package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) dailyplayerstats.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg dailyplayerstats.FetchConfig) ([]*models.DailyPlayerStats, int, error) {
	var err error
	data := []*models.DailyPlayerStats{}
	total := 0
	query := repo.WithParam("SERVER", pg.Safe(cfg.Server)).Model(&data).Context(ctx)

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter).
			Limit(cfg.Filter.Limit).
			Offset(cfg.Filter.Offset)

		order := []string{}

		if cfg.Filter.Sort != "" {
			order = append(order, cfg.Filter.Sort)
		}

		if cfg.Filter.PlayerFilter != nil {
			query = query.Relation("Player._").WhereStruct(cfg.Filter.PlayerFilter)

			if cfg.Filter.PlayerFilter.Sort != "" {
				order = append(order, fmt.Sprintf("player.%s", cfg.Filter.PlayerFilter.Sort))
			}

			if cfg.Filter.PlayerFilter.TribeFilter != nil {
				query = query.
					Join("LEFT JOIN ?SERVER.tribes AS tribe ON tribe.id = player.tribe_id").
					WhereStruct(cfg.Filter.PlayerFilter.TribeFilter)

				if cfg.Filter.PlayerFilter.TribeFilter.Sort != "" {
					order = append(order, fmt.Sprintf("tribe.%s", cfg.Filter.PlayerFilter.TribeFilter.Sort))
				}
			}
		}

		query = query.Order(order...)
	}

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
