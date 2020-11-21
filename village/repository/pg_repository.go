package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) village.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg village.FetchConfig) ([]*models.Village, int, error) {
	var err error
	data := []*models.Village{}
	query := repo.
		WithParam("SERVER", pg.Safe(cfg.Server)).
		Model(&data).
		Context(ctx).
		Order(cfg.Sort...).
		Limit(cfg.Limit).
		Offset(cfg.Offset)
	playerRequired := utils.FindStringWithPrefix(cfg.Sort, "player.") != ""
	tribeRequired := utils.FindStringWithPrefix(cfg.Sort, "tribe.") != ""
	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter)

		if cfg.Filter.XGTE != 0 {
			query = query.Where("x >= ?", cfg.Filter.XGTE)
		} else if cfg.Filter.XGT != 0 {
			query = query.Where("x > ?", cfg.Filter.XGT)
		}
		if cfg.Filter.XLTE != 0 {
			query = query.Where("x <= ?", cfg.Filter.XLTE)
		} else if cfg.Filter.XLT != 0 {
			query = query.Where("x < ?", cfg.Filter.XLT)
		}

		if cfg.Filter.YGTE != 0 {
			query = query.Where("y >= ?", cfg.Filter.YGTE)
		} else if cfg.Filter.YGT != 0 {
			query = query.Where("y > ?", cfg.Filter.YGT)
		}
		if cfg.Filter.YLTE != 0 {
			query = query.Where("y <= ?", cfg.Filter.YLTE)
		} else if cfg.Filter.YLT != 0 {
			query = query.Where("y < ?", cfg.Filter.YLT)
		}

		if len(cfg.Filter.XY) > 0 {
			query = query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				for _, xy := range cfg.Filter.XY {
					splitted := strings.Split(xy, "|")
					if len(splitted) != 2 {
						continue
					}
					x, err := strconv.Atoi(splitted[0])
					if err != nil {
						continue
					}
					y, err := strconv.Atoi(splitted[1])
					if err != nil {
						continue
					}
					q = q.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.Where("x = ?", x)
						q = q.Where("y = ?", y)
						return q, nil
					})
				}
				return q, nil
			})
		}
		if cfg.Filter.PlayerFilter != nil {
			playerRequired = true
			query = query.WhereStruct(cfg.Filter.PlayerFilter)

			if cfg.Filter.PlayerFilter.TribeFilter != nil {
				tribeRequired = true
				query = query.
					WhereStruct(cfg.Filter.PlayerFilter.TribeFilter)
			}
		}

	}
	if playerRequired {
		query = query.Relation("Player._")
	}
	if tribeRequired {
		query = query.Join("LEFT JOIN ?SERVER.tribes AS tribe ON tribe.id = player.tribe_id")
	}
	if len(cfg.Columns) > 0 {
		query = query.Column(cfg.Columns...)
	}

	total := 0
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
