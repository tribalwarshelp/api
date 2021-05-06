package player

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type FetchConfig struct {
	Server string
	Filter *twmodel.PlayerFilter
	Select bool
	Count  bool
	Sort   []string
	Limit  int
	Offset int
}

type SearchPlayerConfig struct {
	Version string
	Name    string
	ID      int
	Count   bool
	Sort    []string
	Limit   int
	Offset  int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Player, int, error)
	FetchNameChanges(ctx context.Context, code twmodel.VersionCode, playerID ...int) (map[int][]*twmodel.PlayerNameChange, error)
	FetchPlayerServers(ctx context.Context, code twmodel.VersionCode, playerID ...int) (map[int][]string, error)
	SearchPlayer(ctx context.Context, cfg SearchPlayerConfig) ([]*twmodel.FoundPlayer, int, error)
}
