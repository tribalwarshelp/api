package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/tribalwarshelp/api/liveennoblement"

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

func NewPGRepository(db *pg.DB, cache redis.UniversalClient) liveennoblement.Repository {
	return &pgRepository{db, cache}
}

func (repo *pgRepository) Fetch(ctx context.Context, server string) ([]*models.LiveEnnoblement, error) {
	if liveennoblements, loaded := repo.loadLiveEnnoblementsFromCache(server); loaded {
		return liveennoblements, nil
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
		fmt.Sprintf(liveennoblement.EndpointGetConquer, time.Now().Add(-1*time.Hour).Unix())
	lines, err := getCSVData(url)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot fetch ennoblements")
	}

	e := []*models.LiveEnnoblement{}
	lineParser := newLineParser()
	for _, line := range lines {
		ennoblement, err := lineParser.parse(line)
		if err != nil {
			continue
		}
		e = append(e, ennoblement)
	}

	go repo.cacheLiveEnnoblements(server, e)

	return e, nil
}

func (repo *pgRepository) loadLiveEnnoblementsFromCache(server string) ([]*models.LiveEnnoblement, bool) {
	ennoblementsJSON, err := repo.cache.Get(context.Background(), fmt.Sprintf(cacheKey, server)).Result()
	if err != nil || ennoblementsJSON == "" {
		return nil, false
	}
	ennoblements := []*models.LiveEnnoblement{}
	if json.Unmarshal([]byte(ennoblementsJSON), &ennoblements) != nil {
		return nil, false
	}
	return ennoblements, true
}

func (repo *pgRepository) cacheLiveEnnoblements(server string, ennoblements []*models.LiveEnnoblement) error {
	ennoblementsJSON, err := json.Marshal(&ennoblements)
	if err != nil {
		return errors.Wrap(err, "cacheLiveEnnoblements")
	}
	if err := repo.cache.Set(context.Background(), fmt.Sprintf(cacheKey, server), ennoblementsJSON, expiration).Err(); err != nil {
		return errors.Wrap(err, "cacheLiveEnnoblements")
	}
	return nil
}
