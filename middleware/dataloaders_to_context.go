package middleware

import (
	"context"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/tribalwarshelp/api/graphql/dataloaders"
	"github.com/tribalwarshelp/api/server"

	"github.com/gin-gonic/gin"
)

var dataloadersContextKey ContextKey = "dataloaders"

func DataLoadersToContext(serverRepo server.Repository, cfg dataloaders.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		loaders := make(map[string]*dataloaders.DataLoaders)
		servers, _, err := serverRepo.Fetch(c.Request.Context(), nil)
		if err != nil {
			c.JSON(http.StatusOK, &gqlerror.Error{
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		for _, server := range servers {
			loaders[server.Key] = dataloaders.New(server.Key, cfg)
		}
		c.Request = c.Request.WithContext(StoreDataLoadersInContext(ctx, loaders))
		c.Next()
	}
}

func StoreDataLoadersInContext(ctx context.Context, loaders map[string]*dataloaders.DataLoaders) context.Context {
	return context.WithValue(ctx, dataloadersContextKey, loaders)
}

func DataLoadersFromContext(ctx context.Context) map[string]*dataloaders.DataLoaders {
	dl := ctx.Value(dataloadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(map[string]*dataloaders.DataLoaders)
}
