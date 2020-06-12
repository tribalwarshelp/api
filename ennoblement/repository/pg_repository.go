package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/tribalwarshelp/api/ennoblement"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

var (
	cacheKey   = "ennoblements-%s"
	expiration = time.Second * 15
)

type pgRepository struct {
	*pg.DB
	cache redis.UniversalClient
}

func NewPGRepository(db *pg.DB, cache redis.UniversalClient) ennoblement.Repository {
	return &pgRepository{db, cache}
}

func (repo *pgRepository) Fetch(ctx context.Context, server string) ([]*models.Ennoblement, error) {
	if ennoblements, loaded := repo.loadEnnoblementsFromCache(server); loaded {
		return ennoblements, nil
	}

	s := &models.Server{}
	if err := repo.Model(s).Where("key = ?", server).Relation("LangVersion").Select(); err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("Server not found")
		}

		return nil, errors.Wrap(err, "Internal server error")
	}

	if s.Status == models.ServerStatusClosed {
		return nil, fmt.Errorf("Server is " + models.ServerStatusClosed.String())
	}

	url := "https://" + s.Key + "." + s.LangVersion.Host +
		fmt.Sprintf(ennoblement.EndpointGetConquer, time.Now().Add(-1*time.Hour).Unix())
	lines, err := getCSVData(url)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot fetch ennoblements")
	}

	e := []*models.Ennoblement{}
	lineParser, err := newLineParser(s.LangVersion.Timezone)
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		ennoblement, err := lineParser.parse(line)
		if err != nil {
			continue
		}
		e = append(e, ennoblement)
	}

	go repo.cacheEnnoblements(server, e)

	return e, nil
}

func (repo *pgRepository) loadEnnoblementsFromCache(server string) ([]*models.Ennoblement, bool) {
	ennoblementsJSON, err := repo.cache.Get(context.Background(), fmt.Sprintf(cacheKey, server)).Result()
	if err != nil || ennoblementsJSON == "" {
		return nil, false
	}
	ennoblements := []*models.Ennoblement{}
	if json.Unmarshal([]byte(ennoblementsJSON), &ennoblements) != nil {
		return nil, false
	}
	return ennoblements, true
}

func (repo *pgRepository) cacheEnnoblements(server string, ennoblements []*models.Ennoblement) error {
	ennoblementsJSON, err := json.Marshal(&ennoblements)
	if err != nil {
		return errors.Wrap(err, "cacheEnnoblements")
	}
	if err := repo.cache.Set(context.Background(), fmt.Sprintf(cacheKey, server), ennoblementsJSON, expiration).Err(); err != nil {
		return errors.Wrap(err, "cacheEnnoblements")
	}
	return nil
}
