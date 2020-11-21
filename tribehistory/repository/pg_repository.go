package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/tribehistory"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) tribehistory.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg tribehistory.FetchConfig) ([]*models.TribeHistory, int, error) {
	var err error
	total := 0
	data := []*models.TribeHistory{}
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Order(cfg.Sort...).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	tribeRequired := utils.FindStringWithPrefix(cfg.Sort, "tribe.") != ""
	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter)
	}
	if tribeRequired {
		query = query.Relation("Tribe._")
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
