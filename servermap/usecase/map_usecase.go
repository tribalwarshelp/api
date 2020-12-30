package usecase

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/tribalwarshelp/api/servermap"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/map-generator/generator"
	"github.com/tribalwarshelp/shared/models"
	"golang.org/x/sync/errgroup"
)

const (
	defaultBarbarianVillageColor = "#808080"
	defaultPlayerVillageColor    = "#FF0000"
)

type usecase struct {
	villageRepo village.Repository
}

func New(villageRepo village.Repository) servermap.Usecase {
	return &usecase{villageRepo}
}

func (ucase *usecase) GetMarkers(ctx context.Context, cfg servermap.GetMarkersConfig) ([]*generator.Marker, error) {
	g := new(errgroup.Group)

	tribes := make(map[string][]int)
	tribeIDs := []int{}
	cache := make(map[int]bool)
	for _, data := range cfg.Tribes {
		//id,#color
		id, color, err := parseMarker(data)
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
		id, color, err := parseMarker(data)
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

	otherMarkers := []*generator.Marker{}
	var otherMarkersMutex sync.Mutex
	if cfg.ShowOtherPlayerVillages {
		color := cfg.PlayerVillageColor
		if color == "" {
			color = defaultPlayerVillageColor
		}
		g.Go(func() error {
			villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
				Server: cfg.Server,
				Filter: &models.VillageFilter{
					PlayerFilter: &models.PlayerFilter{
						IdNEQ:      append(playerIDs, 0),
						TribeIdNEQ: tribeIDs,
					},
				},
				Columns: []string{"x", "y"},
				Count:   false,
			})
			if err != nil {
				return err
			}
			otherMarkersMutex.Lock()
			otherMarkers = append(otherMarkers, &generator.Marker{
				Villages: villages,
				Color:    color,
			})
			otherMarkersMutex.Unlock()
			return nil
		})
	}
	if cfg.ShowBarbarianVillages {
		color := cfg.BarbarianVillageColor
		if color == "" {
			color = defaultBarbarianVillageColor
		}
		g.Go(func() error {
			villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
				Server: cfg.Server,
				Filter: &models.VillageFilter{
					PlayerID: []int{0},
				},
				Columns: []string{"x", "y"},
				Count:   false,
			})
			if err != nil {
				return err
			}
			otherMarkersMutex.Lock()
			otherMarkers = append(otherMarkers, &generator.Marker{
				Villages: villages,
				Color:    color,
			})
			otherMarkersMutex.Unlock()
			return nil
		})
	}

	tribeMarkers := []*generator.Marker{}
	var tribeMarkersMutex sync.Mutex
	for color, tribeIDs := range tribes {
		c := color
		ids := tribeIDs
		g.Go(func() error {
			villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
				Server: cfg.Server,
				Filter: &models.VillageFilter{
					PlayerFilter: &models.PlayerFilter{
						IdNEQ:   playerIDs,
						TribeID: ids,
					},
				},
				Columns: []string{"x", "y"},
				Count:   false,
			})
			if err != nil {
				return err
			}
			tribeMarkersMutex.Lock()
			tribeMarkers = append(tribeMarkers, &generator.Marker{
				Villages: villages,
				Color:    c,
				Larger:   cfg.LargerMarkers,
			})
			tribeMarkersMutex.Unlock()
			return nil
		})
	}

	playerMarkers := []*generator.Marker{}
	var playerMarkersMutex sync.Mutex
	for color, playerIDs := range players {
		c := color
		ids := playerIDs
		g.Go(func() error {
			villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
				Server: cfg.Server,
				Filter: &models.VillageFilter{
					PlayerID: ids,
				},
				Columns: []string{"x", "y"},
				Count:   false,
			})
			if err != nil {
				return err
			}
			playerMarkersMutex.Lock()
			playerMarkers = append(playerMarkers, &generator.Marker{
				Villages: villages,
				Color:    c,
				Larger:   cfg.LargerMarkers,
			})
			playerMarkersMutex.Unlock()
			return nil
		})
	}

	err := g.Wait()
	sort.SliceStable(playerMarkers, func(i, j int) bool {
		return len(playerMarkers[i].Villages) < len(playerMarkers[j].Villages)
	})
	sort.SliceStable(tribeMarkers, func(i, j int) bool {
		return len(tribeMarkers[i].Villages) < len(tribeMarkers[j].Villages)
	})
	return concatMarkers(otherMarkers, tribeMarkers, playerMarkers), err
}

func concatMarkers(slices ...[]*generator.Marker) []*generator.Marker {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]*generator.Marker, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func parseMarker(str string) (int, string, error) {
	splitted := strings.Split(str, ",")
	if len(splitted) != 2 {
		return 0, "", fmt.Errorf("%s: Invalid marker format (should be id,#hexcolor)", str)
	}
	id, err := strconv.Atoi(splitted[0])
	if err != nil {
		return 0, "", errors.Wrapf(err, "%s: Invalid marker format (should be id,#hexcolor)", str)
	}
	if id <= 0 {
		return 0, "", fmt.Errorf("ID should be greater than 0")
	}

	return id, splitted[1], nil
}
