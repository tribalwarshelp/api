package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/utils"
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
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Order(cfg.Sort...).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	playerRequired := utils.FindStringWithPrefix(cfg.Sort, "player") != ""
	tribeRequired := utils.FindStringWithPrefix(cfg.Sort, "tribe") != ""

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter)

		if cfg.Filter.PlayerFilter != nil {
			playerRequired = true
			query = query.WhereStruct(cfg.Filter.PlayerFilter)

			if cfg.Filter.PlayerFilter.TribeFilter != nil {
				tribeRequired = true
				query = query.
					WhereStruct(cfg.Filter.PlayerFilter.TribeFilter)
			}
		}
	}
	if playerRequired {
		query = query.Relation("Player._")
	}
	if tribeRequired {
		query = query.Join("LEFT JOIN ?SERVER.tribes AS tribe ON tribe.id = player.tribe_id")
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
