package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) village.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, server string, f *models.VillageFilter) ([]*models.Village, int, error) {
	var err error
	data := []*models.Village{}
	query := repo.WithParam("SERVER", pg.Safe(server)).Model(&data).Context(ctx)

	if f != nil {
		query = query.
			WhereStruct(f).
			Limit(f.Limit).
			Offset(f.Offset)

		if f.XGTE != 0 {
			query = query.Where("x >= ?", f.XGTE)
		} else if f.XGT != 0 {
			query = query.Where("x > ?", f.XGT)
		}
		if f.XLTE != 0 {
			query = query.Where("x <= ?", f.XLTE)
		} else if f.XLT != 0 {
			query = query.Where("x < ?", f.XLT)
		}

		if f.YGTE != 0 {
			query = query.Where("y >= ?", f.YGTE)
		} else if f.YGT != 0 {
			query = query.Where("y > ?", f.YGT)
		}
		if f.YLTE != 0 {
			query = query.Where("y <= ?", f.YLTE)
		} else if f.YLT != 0 {
			query = query.Where("y < ?", f.YLT)
		}

		if f.Sort != "" {
			query = query.Order(f.Sort)
		}
	}

	total, err := query.SelectAndCount()
	if err != nil && err != pg.ErrNoRows {
		if strings.Contains(err.Error(), `relation "`+server) {
			return nil, 0, fmt.Errorf("Server not found")
		}
		return nil, 0, errors.Wrap(err, "Internal server error")
	}

	return data, total, nil
}
