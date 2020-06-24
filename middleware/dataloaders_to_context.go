package middleware

import (
	"context"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/tribalwarshelp/api/graphql/dataloaders"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/shared/models"

	"github.com/gin-gonic/gin"
)

var serverDataLoadersContextKey ContextKey = "serverDataLoaders"
var langVersionLoadersContextKey ContextKey = "langVersionLoaders"
var dataloadersContextKey ContextKey = "dataloaders"

func DataLoadersToContext(serverRepo server.Repository, cfg dataloaders.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		serverDataLoaders := make(map[string]*dataloaders.ServerDataLoaders)
		langVersionDataLoaders := make(map[models.LanguageTag]*dataloaders.LangVersionDataLoaders)
		servers, _, err := serverRepo.Fetch(c.Request.Context(), server.FetchConfig{})
		if err != nil {
			c.JSON(http.StatusOK, &gqlerror.Error{
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		for _, server := range servers {
			serverDataLoaders[server.Key] = dataloaders.NewServerDataLoaders(server.Key, cfg)
			if _, ok := langVersionDataLoaders[server.LangVersionTag]; !ok {
				langVersionDataLoaders[server.LangVersionTag] = dataloaders.NewLangVersionDataLoaders(server.LangVersionTag, cfg)
			}
		}
		ctx = StoreServerDataLoadersInContext(ctx, serverDataLoaders)
		ctx = StoreLangVersionDataLoadersInContext(ctx, langVersionDataLoaders)
		ctx = StoreDataLoadersInContext(ctx, dataloaders.NewDataLoaders(cfg))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func StoreServerDataLoadersInContext(ctx context.Context, loaders map[string]*dataloaders.ServerDataLoaders) context.Context {
	return context.WithValue(ctx, serverDataLoadersContextKey, loaders)
}

func ServerDataLoadersFromContext(ctx context.Context) map[string]*dataloaders.ServerDataLoaders {
	dl := ctx.Value(serverDataLoadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(map[string]*dataloaders.ServerDataLoaders)
}

func StoreLangVersionDataLoadersInContext(ctx context.Context, loaders map[models.LanguageTag]*dataloaders.LangVersionDataLoaders) context.Context {
	return context.WithValue(ctx, langVersionLoadersContextKey, loaders)
}

func LangVersionDataLoadersFromContext(ctx context.Context) map[models.LanguageTag]*dataloaders.LangVersionDataLoaders {
	dl := ctx.Value(langVersionLoadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(map[models.LanguageTag]*dataloaders.LangVersionDataLoaders)
}

func StoreDataLoadersInContext(ctx context.Context, loaders *dataloaders.DataLoaders) context.Context {
	return context.WithValue(ctx, dataloadersContextKey, loaders)
}

func DataLoadersFromContext(ctx context.Context) *dataloaders.DataLoaders {
	dl := ctx.Value(dataloadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(*dataloaders.DataLoaders)
}
