package resolvers

import (
	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/langversion"
	"github.com/tribalwarshelp/api/liveennoblement"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/playerhistory"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/serverstats"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/tribehistory"
	"github.com/tribalwarshelp/api/village"
)

type Resolver struct {
	LangVersionUcase     langversion.Usecase
	ServerUcase          server.Usecase
	PlayerUcase          player.Usecase
	TribeUcase           tribe.Usecase
	VillageUcase         village.Usecase
	LiveEnnoblementUcase liveennoblement.Usecase
	EnnoblementUcase     ennoblement.Usecase
	PlayerHistoryUcase   playerhistory.Usecase
	TribeHistoryUcase    tribehistory.Usecase
	ServerStatsUcase     serverstats.Usecase
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver     { return &queryResolver{r} }
func (r *Resolver) Player() generated.PlayerResolver   { return &playerResolver{r} }
func (r *Resolver) Village() generated.VillageResolver { return &villageResolver{r} }
func (r *Resolver) LiveEnnoblement() generated.LiveEnnoblementResolver {
	return &liveEnnoblementResolver{r}
}
func (r *Resolver) Ennoblement() generated.EnnoblementResolver { return &ennoblementResolver{r} }
func (r *Resolver) Server() generated.ServerResolver           { return &serverResolver{r} }
func (r *Resolver) ServerStatsRecord() generated.ServerStatsRecordResolver {
	return &serverStatsRecordResolver{r}
}
func (r *Resolver) PlayerHistoryRecord() generated.PlayerHistoryRecordResolver {
	return &playerHistoryRecordResolver{r}
}
func (r *Resolver) TribeHistoryRecord() generated.TribeHistoryRecordResolver {
	return &tribeHistoryRecordResolver{r}
}

type queryResolver struct{ *Resolver }
type playerResolver struct{ *Resolver }
type villageResolver struct{ *Resolver }
type liveEnnoblementResolver struct{ *Resolver }
type ennoblementResolver struct{ *Resolver }
type serverResolver struct{ *Resolver }
type playerHistoryRecordResolver struct{ *Resolver }
type tribeHistoryRecordResolver struct{ *Resolver }
type serverStatsRecordResolver struct{ *Resolver }
