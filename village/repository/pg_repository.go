package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
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
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	if len(cfg.Columns) > 0 {
		query = query.Column(cfg.Columns...)
	}
	relationshipAndSortAppender := &models.VillageRelationshipAndSortAppender{
		Filter: &models.VillageFilter{},
		Sort:   cfg.Sort,
	}
	if cfg.Filter != nil {
		query = query.Apply(cfg.Filter.Where)
		relationshipAndSortAppender.Filter = cfg.Filter
	}
	query = query.Apply(relationshipAndSortAppender.Append)

	total := 0
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
