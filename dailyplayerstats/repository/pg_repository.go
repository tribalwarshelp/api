package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
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
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	relationshipAndSortAppender := &models.DailyPlayerStatsRelationshipAndSortAppender{
		Filter: &models.DailyPlayerStatsFilter{},
		Sort:   cfg.Sort,
	}
	if cfg.Filter != nil {
		query = query.Apply(cfg.Filter.Where)
		relationshipAndSortAppender.Filter = cfg.Filter
	}
	query = query.Apply(relationshipAndSortAppender.Append)

	if cfg.Count {
		total, err = query.SelectAndCount()
	} else {
		err = query.Select()
	}
	if err != nil && err != pg.ErrNoRows {
		if strings.Contains(err.Error(), `relation "`+cfg.Server) {
			return nil, 0, fmt.Errorf("Server not found")
		}
		return nil, 0, fmt.Errorf("Internal server error")
	}

	return data, total, nil
}
