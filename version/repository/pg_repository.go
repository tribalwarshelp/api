package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"

	"github.com/tribalwarshelp/api/version"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) (version.Repository, error) {
	if err := db.Model(&twmodel.Version{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return nil, errors.Wrap(err, "cannot create 'versions' table")
	}
	return &pgRepository{db}, nil
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg version.FetchConfig) ([]*twmodel.Version, int, error) {
	var err error
	data := []*twmodel.Version{}
	total := 0
	query := repo.
		Model(&data).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset).
		Apply(cfg.Filter.Where).
		Apply(gopgutil.OrderAppender{
			Orders:   cfg.Sort,
			MaxDepth: 4,
		}.Apply)

	if cfg.Count && cfg.Select {
		total, err = query.SelectAndCount()
	} else if cfg.Select {
		err = query.Select()
	} else if cfg.Count {
		total, err = query.Count()
	}
	if err != nil && err != pg.ErrNoRows {
		return nil, 0, errors.New("Internal server error")
	}

	return data, total, nil
}
