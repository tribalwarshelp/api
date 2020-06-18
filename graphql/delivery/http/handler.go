package httpdelivery

import (
	"fmt"
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
	playgroundTTL = time.Hour / time.Second
)

type Config struct {
	RouterGroup *gin.RouterGroup
	Resolver    *resolvers.Resolver
}

func Attach(cfg Config) error {
	if cfg.Resolver == nil {
		return fmt.Errorf("Graphql resolver cannot be nil")
	}
	gqlHandler := graphqlHandler(cfg.Resolver)
	cfg.RouterGroup.GET("/graphql", gqlHandler)
	cfg.RouterGroup.POST("/graphql", gqlHandler)
	cfg.RouterGroup.GET("/", playgroundHandler())
	return nil
}

// Defining the GraphQL handler
func graphqlHandler(r *resolvers.Resolver) gin.HandlerFunc {
	cfg := generated.Config{Resolvers: r}
	srv := handler.New(generated.NewExecutableSchema(cfg))

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, must-revalidate")
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("Playground", "/graphql")

	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf(`public, must-revalidate, max-age=%d, s-maxage=%d`, playgroundTTL, playgroundTTL))
		h.ServeHTTP(c.Writer, c.Request)
	}
}
