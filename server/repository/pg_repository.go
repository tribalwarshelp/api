package repository

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) (server.Repository, error) {
	if err := db.CreateTable((*models.Server)(nil), &orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return nil, errors.Wrap(err, "Cannot create 'servers' table")
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
	versionRequired := utils.FindStringWithPrefix(cfg.Sort, "version.") != ""
	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter)
	}
	if versionRequired {
		query = query.Relation("Version._")
	}
	if len(cfg.Columns) > 0 {
		query = query.Column(cfg.Columns...)
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
