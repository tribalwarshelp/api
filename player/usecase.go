package player

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Player, int, error)
	GetByID(ctx context.Context, server string, id int) (*twmodel.Player, error)
	SearchPlayer(ctx context.Context, cfg SearchPlayerConfig) ([]*twmodel.FoundPlayer, int, error)
}
