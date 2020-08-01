package servermap

import (
	"context"

	"github.com/tribalwarshelp/map-generator/generator"
)

type GetMarkersConfig struct {
	Server                  string
	Tribes                  []string
	Players                 []string
	ShowBarbarianVillages   bool
	ShowOtherPlayerVillages bool
	LargerMarkers           bool
}

type Usecase interface {
	GetMarkers(ctx context.Context, cfg GetMarkersConfig) ([]*generator.Marker, error)
}
