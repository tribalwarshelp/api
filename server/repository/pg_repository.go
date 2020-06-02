package repository

import (
	"context"

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
	if err := db.CreateTable((*models.Server)(nil), &orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return nil, errors.Wrap(err, "Cannot create 'servers' table")
	}
	return &pgRepository{db}, nil
}

func (repo *pgRepository) Fetch(ctx context.Context, f *models.ServerFilter) ([]*models.Server, int, error) {
	var err error
	data := []*models.Server{}
	query := repo.Model(&data).Context(ctx)

	if f != nil {
		query = query.
			WhereStruct(f).
			Limit(f.Limit).
			Offset(f.Offset)

		if f.Sort != "" {
			query = query.Order(f.Sort)
		}
	}

	total, err := query.SelectAndCount()
	if err != nil && err != pg.ErrNoRows {
		return nil, 0, errors.Wrap(err, "Internal server error")
	}

	return data, total, nil
}
