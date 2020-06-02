package httpdelivery

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/graphql/resolvers"
)

func Attach(g *gin.RouterGroup, r *resolvers.Resolver) error {
	if r == nil {
		return fmt.Errorf("Graphql resolver cannot be nil")
	}
	g.POST("/graphql", graphqlHandler(r))
	g.GET("/", playgroundHandler())
	return nil
}

// Defining the Graphql handler
func graphqlHandler(r *resolvers.Resolver) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	cfg := generated.Config{Resolvers: r}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, must-revalidate")
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("Playground", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
