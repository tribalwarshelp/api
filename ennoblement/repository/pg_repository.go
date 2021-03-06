package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"strings"

	"github.com/go-pg/pg/v10"

	"github.com/tribalwarshelp/api/ennoblement"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) ennoblement.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg ennoblement.FetchConfig) ([]*twmodel.Ennoblement, int, error) {
	var err error
	total := 0
	data := make([]*twmodel.Ennoblement, 0)
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Limit(cfg.Limit).
		Offset(cfg.Offset).
		Apply(cfg.Filter.WhereWithRelations).
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
		if strings.Contains(err.Error(), `relation "`+cfg.Server) {
			return nil, 0, errors.New("Server not found")
		}
		return nil, 0, errors.New("Internal server error")
	}

	return data, total, nil
}
