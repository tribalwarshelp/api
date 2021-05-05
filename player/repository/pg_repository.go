package repository

import (
	"context"
	"github.com/Kichiyaki/gopgutil/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/tribalwarshelp/api/player"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) player.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, cfg player.FetchConfig) ([]*twmodel.Player, int, error) {
	var err error
	data := []*twmodel.Player{}
	total := 0
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

type fetchPlayerServersQueryResult struct {
	PlayerID int
	Servers  []string `pg:",array"`
}

func (repo *pgRepository) FetchNameChanges(ctx context.Context, code twmodel.VersionCode, playerID ...int) (map[int][]*twmodel.PlayerNameChange, error) {
	data := []*twmodel.PlayerNameChange{}
	if err := repo.Model(&data).
		Context(ctx).
		Where("version_code = ?", code).
		Where("player_id IN (?)", pg.In(playerID)).
		Order("change_date ASC").
		Select(); err != nil && err != pg.ErrNoRows {
		return nil, errors.New("Internal server error")
	}

	m := make(map[int][]*twmodel.PlayerNameChange)
	for _, res := range data {
		m[res.PlayerID] = append(m[res.PlayerID], res)
	}
	return m, nil
}

func (repo *pgRepository) FetchPlayerServers(ctx context.Context, code twmodel.VersionCode, playerID ...int) (map[int][]string, error) {
	data := []*fetchPlayerServersQueryResult{}
	if err := repo.Model(&twmodel.PlayerToServer{}).
		Context(ctx).
		Column("player_id").
		ColumnExpr("array_agg(server_key) as servers").
		Relation("Server._").
		Where("version_code = ?", code).
		Where("player_id IN (?)", pg.In(playerID)).
		Group("player_id").
		Select(&data); err != nil && err != pg.ErrNoRows {
		return nil, errors.New("Internal server error")
	}

	m := make(map[int][]string)
	for _, res := range data {
		m[res.PlayerID] = res.Servers
	}
	return m, nil
}

func (repo *pgRepository) SearchPlayer(ctx context.Context, cfg player.SearchPlayerConfig) ([]*twmodel.FoundPlayer, int, error) {
	servers := []*twmodel.Server{}
	if err := repo.
		Model(&servers).
		Context(ctx).
		Column("key").
		Where("version_code = ?", cfg.Version).
		Select(); err != nil {
		return nil, 0, errors.New("Internal server error")
	}

	var query *orm.Query
	res := []*twmodel.FoundPlayer{}
	whereClause := "player.id = ?1 OR player.name ILIKE ?0"
	if cfg.ID <= 0 {
		whereClause = "player.name ILIKE ?0"
	} else if cfg.Name == "" {
		whereClause = "player.id = ?1"
	}
	for _, server := range servers {
		safeKey := pg.Safe(server.Key)
		otherQuery := repo.
			Model().
			Context(ctx).
			ColumnExpr("? AS server", server.Key).
			ColumnExpr("tribe.tag as tribe_tag").
			Column("player.id", "player.name", "player.most_points", "player.best_rank", "player.most_villages", "player.tribe_id").
			TableExpr("?0.players as player", safeKey).
			Join("LEFT JOIN ?0.tribes as tribe ON player.tribe_id = tribe.id", safeKey).
			Where(whereClause, cfg.Name, cfg.ID)
		if query == nil {
			query = otherQuery
		} else {
			query = query.UnionAll(otherQuery)
		}
	}

	var err error
	count := 0
	if query != nil {
		base := repo.
			Model().
			Context(ctx).
			With("union_q", query).
			Table("union_q").
			Limit(cfg.Limit).
			Offset(cfg.Offset).
			Apply(gopgutil.OrderAppender{
				Orders:   cfg.Sort,
				MaxDepth: 4,
			}.Apply)
		if cfg.Count {
			count, err = base.SelectAndCount(&res)
		} else {
			err = base.Select(&res)
		}
		if err != nil && err != pg.ErrNoRows {
			return nil, 0, errors.New("Internal server error")
		}
	}

	return res, count, nil
}
