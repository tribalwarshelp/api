package querycomplexity

import (
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

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
	countComplexity                      = 1000
	dailyPlayerStatsTotalFieldComplexity = 1000
	dailyTribeStatsTotalFieldComplexity  = dailyPlayerStatsTotalFieldComplexity
	ennoblementsTotalFieldComplexity     = 1000
	playerHistoryTotalFieldComplexity    = 500
	tribeHistoryTotalFieldComplexity     = playersTotalFieldComplexity
	tribeChangesTotalFieldComplexity     = 300
	searchPlayerTotalFieldComplexity     = 1000
	searchTribeTotalFieldComplexity      = searchPlayerTotalFieldComplexity
	playersTotalFieldComplexity          = 300
	tribesTotalFieldComplexity           = playersTotalFieldComplexity
	villagesTotalFieldComplexity         = 1000
	serverStatsTotalFieldComplexity      = 100
	serversTotalFieldComplexity          = 200
	versionsTotalFieldComplexity         = 50
)

func GetComplexityRoot() generated.ComplexityRoot {
	complexityRoot := generated.ComplexityRoot{}

	complexityRoot.DailyPlayerStats.Total = getCountComplexity
	complexityRoot.Query.DailyPlayerStats = func(
		childComplexity int,
		server string,
		filter *twmodel.DailyPlayerStatsFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, dailyplayerstats.FetchLimit),
			dailyPlayerStatsTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.DailyTribeStats.Total = getCountComplexity
	complexityRoot.Query.DailyTribeStats = func(
		childComplexity int,
		server string,
		filter *twmodel.DailyTribeStatsFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, dailytribestats.FetchLimit),
			dailyTribeStatsTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.EnnoblementList.Total = getCountComplexity
	complexityRoot.Query.Ennoblements = func(
		childComplexity int,
		server string,
		filter *twmodel.EnnoblementFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, ennoblement.FetchLimit),
			ennoblementsTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.PlayerHistory.Total = getCountComplexity
	complexityRoot.Query.PlayerHistory = func(
		childComplexity int,
		server string,
		filter *twmodel.PlayerHistoryFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, playerhistory.FetchLimit),
			playerHistoryTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.TribeHistory.Total = getCountComplexity
	complexityRoot.Query.TribeHistory = func(
		childComplexity int,
		server string,
		filter *twmodel.TribeHistoryFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, tribehistory.FetchLimit),
			tribeHistoryTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.TribeChanges.Total = getCountComplexity
	complexityRoot.Query.TribeChanges = func(
		childComplexity int,
		server string,
		filter *twmodel.TribeChangeFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, tribechange.FetchLimit),
			tribeChangesTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.FoundPlayerList.Total = getCountComplexity
	complexityRoot.Query.SearchPlayer = func(
		childComplexity int,
		version string,
		name *string,
		id *int,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, player.FetchLimit),
			searchPlayerTotalFieldComplexity,
			3,
		)
	}

	complexityRoot.FoundTribeList.Total = getCountComplexity
	complexityRoot.Query.SearchTribe = func(
		childComplexity int,
		version string,
		query string,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, tribe.FetchLimit),
			searchTribeTotalFieldComplexity,
			3,
		)
	}

	complexityRoot.Player.NameChanges = func(childComplexity int) int {
		return 10 + childComplexity
	}
	complexityRoot.Player.Servers = func(childComplexity int) int {
		return 10 + childComplexity
	}
	complexityRoot.PlayerList.Total = getCountComplexity
	complexityRoot.Query.Players = func(
		childComplexity int,
		server string,
		filter *twmodel.PlayerFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, player.FetchLimit),
			playersTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.TribeList.Total = getCountComplexity
	complexityRoot.Query.Tribes = func(
		childComplexity int,
		server string,
		filter *twmodel.TribeFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, tribe.FetchLimit),
			tribesTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.VillageList.Total = getCountComplexity
	complexityRoot.Query.Villages = func(
		childComplexity int,
		server string,
		filter *twmodel.VillageFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, village.FetchLimit),
			villagesTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.ServerStats.Total = getCountComplexity
	complexityRoot.Query.ServerStats = func(
		childComplexity int,
		server string,
		filter *twmodel.ServerStatsFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, serverstats.FetchLimit),
			serverStatsTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.ServerList.Total = getCountComplexity
	complexityRoot.Query.Servers = func(
		childComplexity int,
		filter *twmodel.ServerFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, server.FetchLimit),
			serversTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.VersionList.Total = getCountComplexity
	complexityRoot.Query.Versions = func(
		childComplexity int,
		filter *twmodel.VersionFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			safeptr.SafeIntPointer(limit, version.FetchLimit),
			versionsTotalFieldComplexity,
			1,
		)
	}

	return complexityRoot
}

func computeComplexity(childComplexity, limit, totalFieldComplexity, multiplyBy int) int {
	complexity := 0
	if childComplexity >= countComplexity {
		childComplexity -= countComplexity
		complexity += totalFieldComplexity
	}
	return limit*childComplexity*multiplyBy + complexity
}

func getCountComplexity(childComplexity int) int {
	return countComplexity + childComplexity
}
