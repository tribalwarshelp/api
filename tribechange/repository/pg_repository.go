package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/tribechange"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) tribechange.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg tribechange.FetchConfig) ([]*models.TribeChange, int, error) {
	var err error
	total := 0
	data := []*models.TribeChange{}
	query := repo.WithParam("SERVER", pg.Safe(cfg.Server)).Model(&data).Context(ctx)

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter).
			Limit(cfg.Filter.Limit).
			Offset(cfg.Filter.Offset)

		if cfg.Filter.Sort != "" {
			query = query.Order(cfg.Filter.Sort)
		}

		if cfg.Filter.Or != nil {
			query = query.WhereGroup(appendTribeChangeFilterOr(cfg.Filter.Or))
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
