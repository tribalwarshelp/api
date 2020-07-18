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

func (repo *pgRepository) Fetch(ctx context.Context, cfg player.FetchConfig) ([]*models.Player, int, error) {
	var err error
	data := []*models.Player{}
	total := 0
	query := repo.WithParam("SERVER", pg.Safe(cfg.Server)).Model(&data).Context(ctx)

	if cfg.Filter != nil {
		query = query.
			WhereStruct(cfg.Filter).
			Limit(cfg.Filter.Limit).
			Offset(cfg.Filter.Offset)

		if cfg.Filter.Exists != nil {
			query = query.Where("exists = ?", *cfg.Filter.Exists)
		}

		order := []string{}

		if cfg.Filter.Sort != "" {
			order = append(order, cfg.Filter.Sort)
		}

		if cfg.Filter.TribeFilter != nil {
			query = query.Relation("Tribe._").WhereStruct(cfg.Filter.TribeFilter)

			if cfg.Filter.TribeFilter.Sort != "" {
				order = append(order, fmt.Sprintf("tribe.%s", cfg.Filter.TribeFilter.Sort))
			}
		}

		query = query.Order(order...)
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

type fetchPlayerServersQueryResult struct {
	PlayerID int
	Servers  []string `pg:",array"`
}

func (repo *pgRepository) FetchNameChanges(ctx context.Context, langTag models.LanguageTag, playerID ...int) (map[int][]*models.PlayerNameChange, error) {
	data := []*models.PlayerNameChange{}
	if err := repo.Model(&data).
		Context(ctx).
		Where("lang_version_tag = ?", langTag).
		Where("player_id IN (?)", pg.In(playerID)).
		Select(); err != nil && err != pg.ErrNoRows {
		return nil, errors.Wrap(err, "Internal server error")
	}

	m := make(map[int][]*models.PlayerNameChange)
	for _, res := range data {
		m[res.PlayerID] = append(m[res.PlayerID], res)
	}
	return m, nil
}

func (repo *pgRepository) FetchPlayerServers(ctx context.Context, langTag models.LanguageTag, playerID ...int) (map[int][]string, error) {
	data := []*fetchPlayerServersQueryResult{}
	if err := repo.Model(&models.PlayerToServer{}).
		Context(ctx).
		Column("player_id").
		ColumnExpr("array_agg(server_key) as servers").
		Relation("Server._").
		Where("lang_version_tag = ?", langTag).
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
