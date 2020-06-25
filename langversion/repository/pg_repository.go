package repository

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/langversion"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) (langversion.Repository, error) {
	if err := db.CreateTable((*models.LangVersion)(nil), &orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return nil, errors.Wrap(err, "Cannot create 'lang_versions' table")
	}
	return &pgRepository{db}, nil
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg langversion.FetchConfig) ([]*models.LangVersion, int, error) {
	var err error
	data := []*models.LangVersion{}
	total := 0
	query := repo.Model(&data).Context(ctx)

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter).
			Limit(cfg.Filter.Limit).
			Offset(cfg.Filter.Offset)

		if cfg.Filter.Sort != "" {
			query = query.Order(cfg.Filter.Sort)
		}
	}

	if cfg.Count {
		total, err = query.SelectAndCount()
	} else {
		err = query.Select()
	}
	if err != nil && err != pg.ErrNoRows {
		return nil, 0, errors.Wrap(err, "Internal server error")
	}

	return data, total, nil
}
