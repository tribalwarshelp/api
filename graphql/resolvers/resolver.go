package resolvers

import (
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/langversion"
	"github.com/tribalwarshelp/api/liveennoblement"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/village"
)

type Resolver struct {
	LangVersionUcase     langversion.Usecase
	ServerUcase          server.Usecase
	PlayerUcase          player.Usecase
	TribeUcase           tribe.Usecase
	VillageUcase         village.Usecase
	LiveEnnoblementUcase liveennoblement.Usecase
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver     { return &queryResolver{r} }
func (r *Resolver) Player() generated.PlayerResolver   { return &playerResolver{r} }
func (r *Resolver) Village() generated.VillageResolver { return &villageResolver{r} }
func (r *Resolver) LiveEnnoblement() generated.LiveEnnoblementResolver {
	return &liveEnnoblementResolver{r}
}
func (r *Resolver) Server() generated.ServerResolver { return &serverResolver{r} }

type queryResolver struct{ *Resolver }
type playerResolver struct{ *Resolver }
type villageResolver struct{ *Resolver }
type liveEnnoblementResolver struct{ *Resolver }
type serverResolver struct{ *Resolver }
