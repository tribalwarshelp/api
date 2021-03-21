package repository

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) (server.Repository, error) {
	if err := db.Model(&models.Server{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return nil, errors.Wrap(err, "cannot create 'servers' table")
	}
	return &pgRepository{db}, nil
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg server.FetchConfig) ([]*models.Server, int, error) {
	var err error
	total := 0
	data := []*models.Server{}
	query := repo.
		Model(&data).
		Context(ctx).
		Order(cfg.Sort...).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	if cfg.Filter != nil {
		query = query.Apply(cfg.Filter.Where)
	}
	if len(cfg.Columns) > 0 {
		query = query.Column(cfg.Columns...)
	}

	if cfg.Count && cfg.Select {
		total, err = query.SelectAndCount()
	} else if cfg.Select {
		err = query.Select()
	} else if cfg.Count {
		total, err = query.Count()
	}
	if err != nil && err != pg.ErrNoRows {
		return nil, 0, fmt.Errorf("Internal server error")
	}

	return data, total, nil
}
