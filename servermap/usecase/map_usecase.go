package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/servermap"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/map-generator/generator"
	"github.com/tribalwarshelp/shared/models"
)

const (
	defaultBarbarianVillagesColor = "#808080"
	defaultPlayerVillagesColor    = "#FF0000"
)

type usecase struct {
	villageRepo village.Repository
}

func New(villageRepo village.Repository) servermap.Usecase {
	return &usecase{villageRepo}
}

func (ucase *usecase) GetMarkers(ctx context.Context, cfg servermap.GetMarkersConfig) ([]*generator.Marker, error) {
	tribes := make(map[string][]int)
	tribeIDs := []int{}
	cache := make(map[int]bool)
	for _, data := range cfg.Tribes {
		//id,#color
		id, color, err := parseQueryParam(data)
		if err != nil {
			return nil, errors.Wrapf(err, "tribe=%s", data)
		}
		if ok := cache[id]; ok || color == "" {
			continue
		}
		tribeIDs = append(tribeIDs, id)
		cache[id] = true
		tribes[color] = append(tribes[color], id)
	}

	players := make(map[string][]int)
	playerIDs := []int{}
	cache = make(map[int]bool)
	for _, data := range cfg.Players {
		//id,#color
		id, color, err := parseQueryParam(data)
		if err != nil {
			return nil, errors.Wrapf(err, "player=%s", data)
		}
		if ok := cache[id]; ok || color == "" {
			continue
		}
		playerIDs = append(playerIDs, id)
		cache[id] = true
		players[color] = append(players[color], id)
	}

	markers := []*generator.Marker{}

	if cfg.ShowOtherPlayerVillages {
		villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
			Server: cfg.Server,
			Filter: &models.VillageFilter{
				PlayerFilter: &models.PlayerFilter{
					IdNEQ: append(playerIDs, 0),
					TribeFilter: &models.TribeFilter{
						IdNEQ: tribeIDs,
					},
				},
			},
			Count: false,
		})
		if err != nil {
			return nil, err
		}
		markers = append(markers, &generator.Marker{
			Villages: villages,
			Color:    defaultPlayerVillagesColor,
		})
	}

	if cfg.ShowBarbarianVillages {
		villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
			Server: cfg.Server,
			Filter: &models.VillageFilter{
				PlayerID: []int{0},
			},
			Count: false,
		})
		if err != nil {
			return nil, err
		}
		markers = append(markers, &generator.Marker{
			Villages: villages,
			Color:    defaultBarbarianVillagesColor,
		})
	}

	for color, tribeIDs := range tribes {
		villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
			Server: cfg.Server,
			Filter: &models.VillageFilter{
				PlayerFilter: &models.PlayerFilter{
					TribeID: tribeIDs,
				},
			},
			Count: false,
		})
		if err != nil {
			return nil, err
		}
		markers = append(markers, &generator.Marker{
			Villages: villages,
			Color:    color,
		})
	}

	for color, playerIDs := range players {
		villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
			Server: cfg.Server,
			Filter: &models.VillageFilter{
				PlayerID: playerIDs,
			},
			Count: false,
		})
		if err != nil {
			return nil, err
		}
		markers = append(markers, &generator.Marker{
			Villages: villages,
			Color:    color,
		})
	}

	return markers, nil
}

func parseQueryParam(str string) (int, string, error) {
	splitted := strings.Split(str, ",")
	if len(splitted) != 2 {
		return 0, "", fmt.Errorf("Invalid format (should be id,#hexcolor)")
	}
	id, err := strconv.Atoi(splitted[0])
	if err != nil {
		return 0, "", errors.Wrapf(err, "Invalid format (should be id,hexcolor)")
	}
	if id <= 0 {
		return 0, "", fmt.Errorf("ID should be greater than 0")
	}
	return id, splitted[1], nil
}