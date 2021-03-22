package httpdelivery

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/tribalwarshelp/api/graphql/querycomplexity"
	"github.com/tribalwarshelp/api/middleware"
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
		Complexity: querycomplexity.GetComplexityRoot(),
	}
}
