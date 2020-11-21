package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/tribechange"
	"github.com/tribalwarshelp/api/utils"
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
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Order(cfg.Sort...).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	playerRequired := utils.FindStringWithPrefix(cfg.Sort, "player.") != ""
	oldTribeRequired := utils.FindStringWithPrefix(cfg.Sort, "old_tribe.") != ""
	newTribeRequired := utils.FindStringWithPrefix(cfg.Sort, "new_tribe.") != ""
	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter)

		if cfg.Filter.Or != nil {
			query = query.WhereGroup(appendTribeChangeFilterOr(cfg.Filter.Or))
		}
	}
	if playerRequired {
		query = query.Relation("Village._")
	}
	if oldTribeRequired {
		query = query.Join("LEFT JOIN ?SERVER.tribes AS old_tribe ON old_tribe.id = ennoblement.old_tribe_id")
	}
	if newTribeRequired {
		query = query.Join("LEFT JOIN ?SERVER.tribes AS new_tribe ON new_tribe.id = ennoblement.new_tribe_id")
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
