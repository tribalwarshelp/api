package httpdelivery

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/dailytribestats"
	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/middleware"
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
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/graphql/resolvers"
)

const (
	endpointMain       = "/graphql"
	endpointPlayground = "/"
	playgroundTTL      = time.Hour / time.Second
)

type Config struct {
	RouterGroup *gin.RouterGroup
	Resolver    *resolvers.Resolver
}

func Attach(cfg Config) error {
	if cfg.Resolver == nil {
		return fmt.Errorf("Graphql resolver cannot be nil")
	}
	gqlHandler := graphqlHandler(prepareConfig(cfg.Resolver))
	cfg.RouterGroup.GET(endpointMain, gqlHandler)
	cfg.RouterGroup.POST(endpointMain, gqlHandler)
	cfg.RouterGroup.GET(endpointPlayground, playgroundHandler())
	return nil
}

// Defining the GraphQL handler
func graphqlHandler(cfg generated.Config) gin.HandlerFunc {
	srv := handler.New(generated.NewExecutableSchema(cfg))

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(&extension.ComplexityLimit{
		Func: func(ctx context.Context, rc *graphql.OperationContext) int {
			if middleware.CanExceedLimit(ctx) {
				return 500000000
			}
			return 18000
		},
	})

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, must-revalidate")
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("Playground", endpointMain)

	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf(`public, max-age=%d`, playgroundTTL))
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func prepareConfig(r *resolvers.Resolver) generated.Config {
	return generated.Config{
		Resolvers:  r,
		Complexity: getComplexityRoot(),
	}
}

func getComplexityRoot() generated.ComplexityRoot {
	complexityRoot := generated.ComplexityRoot{}
	complexityRoot.Player.NameChanges = func(childComplexity int) int {
		return 10 + childComplexity
	}
	complexityRoot.Player.Servers = func(childComplexity int) int {
		return 10 + childComplexity
	}
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

func computeComplexity(childComplexity, limit, defaultLimit, multiplyBy int) int {
	if childComplexity <= 1 {
		return defaultLimit / 2
	}
	return limit * childComplexity * multiplyBy
}
