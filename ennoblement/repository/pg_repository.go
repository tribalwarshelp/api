package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) ennoblement.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg ennoblement.FetchConfig) ([]*models.Ennoblement, int, error) {
	var err error
	total := 0
	data := []*models.Ennoblement{}
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Order(cfg.Sort...).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	villageRequired := utils.FindStringWithPrefix(cfg.Sort, "village.") != ""
	newOwnerRequired := utils.FindStringWithPrefix(cfg.Sort, "new_owner.") != ""
	newOwnerTribeRequired := utils.FindStringWithPrefix(cfg.Sort, "new_owner_tribe.") != ""
	oldOwnerRequired := utils.FindStringWithPrefix(cfg.Sort, "old_owner.") != ""
	oldOwnerTribeRequired := utils.FindStringWithPrefix(cfg.Sort, "old_owner_tribe.") != ""

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter).
			WhereGroup(appendOrFilter(cfg.Filter.Or))
	}

	if villageRequired {
		query = query.Relation("Village._")
	}
	if newOwnerRequired {
		query = query.Join("LEFT JOIN ?SERVER.players AS new_owner ON new_owner.id = ennoblement.new_owner_id")
	}
	if newOwnerTribeRequired {
		query = query.Join("LEFT JOIN ?SERVER.tribes AS new_owner_tribe ON new_owner_tribe.id = ennoblement.new_owner_tribe_id")
	}
	if oldOwnerRequired {
		query = query.Join("LEFT JOIN ?SERVER.players AS old_owner ON old_owner.id = ennoblement.old_owner_id")
	}
	if oldOwnerTribeRequired {
		query = query.Join("LEFT JOIN ?SERVER.tribes AS old_owner_tribe ON old_owner_tribe.id = ennoblement.old_owner_tribe_id")
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
