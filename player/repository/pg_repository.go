package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) player.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, server string, f *models.PlayerFilter) ([]*models.Player, int, error) {
	var err error
	data := []*models.Player{}
	query := repo.WithParam("SERVER", pg.Safe(server)).Model(&data).Context(ctx)

	if f != nil {
		query = query.
			WhereStruct(f).
			Limit(f.Limit).
			Offset(f.Offset)

		if f.Sort != "" {
			query = query.Order(f.Sort)
		}

		if f.Exist != nil {
			query = query.Where("exist = ?", *f.Exist)
		}

		if f.TribeFilter != nil {
			query = query.Relation("Tribe._").WhereStruct(f.TribeFilter)
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

type fetchPlayerServersQueryResult struct {
	PlayerID int
	Servers  []string `pg:",array"`
}

func (repo *pgRepository) FetchPlayerServers(ctx context.Context, playerID ...int) (map[int][]string, error) {
	data := []*fetchPlayerServersQueryResult{}
	if err := repo.Model(&models.PlayerToServer{}).
		Context(ctx).
		Column("player_id").
		ColumnExpr("array_agg(server_key) as servers").
		Where("player_id IN (?)", pg.In(playerID)).
		Group("player_id").
		Select(&data); err != nil && err != pg.ErrNoRows {
		return nil, errors.Wrap(err, "Internal server error")
	}

	m := make(map[int][]string)
	for _, res := range data {
		m[res.PlayerID] = res.Servers
	}
	return m, nil
}
