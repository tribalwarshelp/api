package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
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
	query := repo.WithParam("SERVER", pg.Safe(cfg.Server)).Model(&data).Context(ctx)

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter).
			Limit(cfg.Filter.Limit).
			Offset(cfg.Filter.Offset)

		if cfg.Filter.Sort != "" {
			query = query.Order(cfg.Filter.Sort)
		}

		if cfg.Filter.Exists != nil {
			query = query.Where("exist = ?", *cfg.Filter.Exists)
		}
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
