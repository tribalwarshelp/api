package querycomplexity

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
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/api/version"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

const (
	countComplexity = 1000
)

func GetComplexityRoot() generated.ComplexityRoot {
	complexityRoot := generated.ComplexityRoot{}
	complexityRoot.Player.NameChanges = func(childComplexity int) int {
		return 10 + childComplexity
	}
	complexityRoot.Player.Servers = func(childComplexity int) int {
		return 10 + childComplexity
	}
	complexityRoot.DailyPlayerStats.Total = getCountComplexity
	complexityRoot.Query.DailyPlayerStats = func(
		childComplexity int,
		server string,
		filter *models.DailyPlayerStatsFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, dailyplayerstats.FetchLimit),
			dailyplayerstats.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.DailyTribeStats = func(
		childComplexity int,
		server string,
		filter *models.DailyTribeStatsFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, dailytribestats.FetchLimit),
			dailytribestats.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.Ennoblements = func(
		childComplexity int,
		server string,
		filter *models.EnnoblementFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, ennoblement.FetchLimit),
			ennoblement.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.PlayerHistory = func(
		childComplexity int,
		server string,
		filter *models.PlayerHistoryFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, playerhistory.FetchLimit),
			playerhistory.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.TribeHistory = func(
		childComplexity int,
		server string,
		filter *models.TribeHistoryFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, tribehistory.FetchLimit),
			tribehistory.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.TribeChanges = func(
		childComplexity int,
		server string,
		filter *models.TribeChangeFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, tribechange.FetchLimit),
			tribechange.FetchLimit,
			1,
		)
	}
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
			utils.SafeIntPointer(limit, player.FetchLimit),
			player.FetchLimit,
			3,
		)
	}
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
			utils.SafeIntPointer(limit, tribe.FetchLimit),
			tribe.FetchLimit,
			3,
		)
	}
	complexityRoot.Query.Players = func(
		childComplexity int,
		server string,
		filter *models.PlayerFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, player.FetchLimit),
			player.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.Tribes = func(
		childComplexity int,
		server string,
		filter *models.TribeFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, tribe.FetchLimit),
			tribe.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.Villages = func(
		childComplexity int,
		server string,
		filter *models.VillageFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, village.FetchLimit),
			village.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.ServerStats = func(
		childComplexity int,
		server string,
		filter *models.ServerStatsFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, serverstats.FetchLimit),
			serverstats.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.Servers = func(
		childComplexity int,
		filter *models.ServerFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, server.FetchLimit),
			server.FetchLimit,
			1,
		)
	}
	complexityRoot.Query.Versions = func(
		childComplexity int,
		filter *models.VersionFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, version.FetchLimit),
			version.FetchLimit,
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
