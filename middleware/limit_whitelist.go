package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

var limitWhitelistContextKey ContextKey = "limitWhitelist"

type LimitWhitelistConfig struct {
	IPAddresses []string
}

func LimitWhitelist(cfg LimitWhitelistConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		clientIP := c.ClientIP()
		mayExceedLimit := false
		for _, ip := range cfg.IPAddresses {
			if ip == clientIP {
				mayExceedLimit = true
				break
			}
		}
		ctx = StoreLimitWhitelistDataInContext(ctx, mayExceedLimit)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func StoreLimitWhitelistDataInContext(ctx context.Context, mayExceedLimit bool) context.Context {
	return context.WithValue(ctx, limitWhitelistContextKey, mayExceedLimit)
}

func MayExceedLimit(ctx context.Context) bool {
	whitelisted := ctx.Value(limitWhitelistContextKey)
	if whitelisted == nil {
		return false
	}
	return whitelisted.(bool)
}
