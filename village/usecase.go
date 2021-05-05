package village

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Village, int, error)
	GetByID(ctx context.Context, server string, id int) (*twmodel.Village, error)
}
