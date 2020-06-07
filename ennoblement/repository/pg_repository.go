package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/tribalwarshelp/api/ennoblement"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type pgRepository struct {
	*pg.DB
}

func NewPGRepository(db *pg.DB) ennoblement.Repository {
	return &pgRepository{db}
}

func (repo *pgRepository) Fetch(ctx context.Context, server string) ([]*models.Ennoblement, error) {
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

	return e, nil
}
