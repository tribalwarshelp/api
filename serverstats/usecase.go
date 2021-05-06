package serverstats

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.ServerStats, int, error)
}
