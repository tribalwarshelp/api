package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/tribalwarshelp/map-generator/generator"

	"github.com/tribalwarshelp/api/servermap"
	"github.com/tribalwarshelp/api/village"

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

	tribes, tribeIDs, err := toMarkers(cfg.Tribes)
	if err != nil {
		return nil, err
	}

	players, playerIDs, err := toMarkers(cfg.Players)
	if err != nil {
		return nil, err
	}

	var otherMarkers []*generator.Marker
	var otherMarkersMutex sync.Mutex
	if cfg.ShowOtherPlayerVillages {
		color := cfg.PlayerVillageColor
		if color == "" {
			color = defaultPlayerVillageColor
		}
		g.Go(func() error {
			villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
				Server: cfg.Server,
				Filter: &twmodel.VillageFilter{
					PlayerFilter: &twmodel.PlayerFilter{
						IDNEQ:      append(playerIDs, 0),
						TribeIDNEQ: tribeIDs,
					},
				},
				Select:  true,
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
				Filter: &twmodel.VillageFilter{
					PlayerID: []int{0},
				},
				Select:  true,
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

	var tribeMarkers []*generator.Marker
	var tribeMarkersMutex sync.Mutex
	for color, tribeIDs := range tribes {
		c := color
		ids := tribeIDs
		g.Go(func() error {
			villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
				Server: cfg.Server,
				Filter: &twmodel.VillageFilter{
					PlayerFilter: &twmodel.PlayerFilter{
						IDNEQ:   playerIDs,
						TribeID: ids,
					},
				},
				Select:  true,
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

	var playerMarkers []*generator.Marker
	var playerMarkersMutex sync.Mutex
	for color, playerIDs := range players {
		c := color
		ids := playerIDs
		g.Go(func() error {
			villages, _, err := ucase.villageRepo.Fetch(ctx, village.FetchConfig{
				Server: cfg.Server,
				Filter: &twmodel.VillageFilter{
					PlayerID: ids,
				},
				Select:  true,
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

	err = g.Wait()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(playerMarkers, func(i, j int) bool {
		return len(playerMarkers[i].Villages) < len(playerMarkers[j].Villages)
	})
	sort.SliceStable(tribeMarkers, func(i, j int) bool {
		return len(tribeMarkers[i].Villages) < len(tribeMarkers[j].Villages)
	})
	return concatMarkers(otherMarkers, tribeMarkers, playerMarkers), nil
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

func toMarker(param string) (int, string, error) {
	parts := strings.Split(param, ",")
	if len(parts) != 2 {
		return 0, "", errors.Errorf("%s: Invalid marker format (should be id,#hexcolor)", param)
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", errors.Wrapf(err, "%s: Invalid marker format (should be id,#hexcolor)", param)
	}
	if id <= 0 {
		return 0, "", errors.New("ID should be greater than 0")
	}

	return id, parts[1], nil
}

func toMarkers(params []string) (map[string][]int, []int, error) {
	idsByColor := make(map[string][]int)
	var ids []int
	cache := make(map[int]bool)
	count := 0
	for _, param := range params {
		if count >= servermap.MaxMarkers {
			break
		}
		//id,#color
		id, color, err := toMarker(param)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "invalid param %s", param)
		}
		if ok := cache[id]; ok || color == "" {
			continue
		}
		ids = append(ids, id)
		cache[id] = true
		idsByColor[color] = append(idsByColor[color], id)
		count++
	}
	return idsByColor, ids, nil
}
