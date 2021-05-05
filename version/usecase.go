package version

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Version, int, error)
	GetByCode(ctx context.Context, code twmodel.VersionCode) (*twmodel.Version, error)
}
