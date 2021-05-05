package server

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Server, int, error)
	GetByKey(ctx context.Context, key string) (*twmodel.Server, error)
}
