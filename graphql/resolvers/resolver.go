package resolvers

import (
	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/dailytribestats"
	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/playerhistory"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/serverstats"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/tribechange"
	"github.com/tribalwarshelp/api/tribehistory"
	"github.com/tribalwarshelp/api/version"
	"github.com/tribalwarshelp/api/village"
)

const (
	countField = "total"
	itemsField = "items"
)

type Resolver struct {
	VersionUcase          version.Usecase
	ServerUcase           server.Usecase
	PlayerUcase           player.Usecase
	TribeUcase            tribe.Usecase
	VillageUcase          village.Usecase
	EnnoblementUcase      ennoblement.Usecase
	PlayerHistoryUcase    playerhistory.Usecase
	TribeHistoryUcase     tribehistory.Usecase
	ServerStatsUcase      serverstats.Usecase
	TribeChangeUcase      tribechange.Usecase
	DailyTribeStatsUcase  dailytribestats.Usecase
	DailyPlayerStatsUcase dailyplayerstats.Usecase
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver             { return &queryResolver{r} }
func (r *Resolver) Player() generated.PlayerResolver           { return &playerResolver{r} }
func (r *Resolver) Village() generated.VillageResolver         { return &villageResolver{r} }
func (r *Resolver) Ennoblement() generated.EnnoblementResolver { return &ennoblementResolver{r} }
func (r *Resolver) Server() generated.ServerResolver           { return &serverResolver{r} }
func (r *Resolver) PlayerHistoryRecord() generated.PlayerHistoryRecordResolver {
	return &playerHistoryRecordResolver{r}
}
func (r *Resolver) TribeHistoryRecord() generated.TribeHistoryRecordResolver {
	return &tribeHistoryRecordResolver{r}
}
func (r *Resolver) TribeChangeRecord() generated.TribeChangeRecordResolver {
	return &tribeChangeRecordResolver{r}
}
func (r *Resolver) DailyPlayerStatsRecord() generated.DailyPlayerStatsRecordResolver {
	return &dailyPlayerStatsRecordResolver{r}
}
func (r *Resolver) DailyTribeStatsRecord() generated.DailyTribeStatsRecordResolver {
	return &dailyTribeStatsRecordResolver{r}
}

type queryResolver struct{ *Resolver }
type playerResolver struct{ *Resolver }
type villageResolver struct{ *Resolver }
type tribeResolver struct{ *Resolver }
type ennoblementResolver struct{ *Resolver }
type serverResolver struct{ *Resolver }
type playerHistoryRecordResolver struct{ *Resolver }
type tribeHistoryRecordResolver struct{ *Resolver }
type serverStatsRecordResolver struct{ *Resolver }
type tribeChangeRecordResolver struct{ *Resolver }
type dailyPlayerStatsRecordResolver struct{ *Resolver }
type dailyTribeStatsRecordResolver struct{ *Resolver }
type versionResolver struct{ *Resolver }
