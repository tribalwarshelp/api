package middleware

import (
	"context"
	"github.com/Kichiyaki/appmode"
	"net"

	"github.com/gin-gonic/gin"
)

var limitWhitelistContextKey ContextKey = "limitWhitelist"

type networksAndIPs struct {
	Networks []*net.IPNet
	IPs      []net.IP
}

func (networksAndIps networksAndIPs) Contains(ip net.IP) bool {
	for _, whitelistedIP := range networksAndIps.IPs {
		if whitelistedIP.Equal(ip) {
			return true
		}
	}
	for _, network := range networksAndIps.Networks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

type LimitWhitelistConfig struct {
	IPAddresses []string
}

func (cfg LimitWhitelistConfig) getNetworksAndIps() networksAndIPs {
	var networks []*net.IPNet
	var ips []net.IP
	for _, ip := range cfg.IPAddresses {
		_, network, err := net.ParseCIDR(ip)
		if err == nil {
			networks = append(networks, network)
			continue
		}

		parsed := net.ParseIP(ip)
		if parsed != nil {
			ips = append(ips, parsed)
		}
	}
	return networksAndIPs{
		Networks: networks,
		IPs:      ips,
	}
}

func LimitWhitelist(cfg LimitWhitelistConfig) gin.HandlerFunc {
	networksAndIps := cfg.getNetworksAndIps()
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		clientIP := net.ParseIP(c.ClientIP())
		canExceedLimit := networksAndIps.Contains(clientIP) || appmode.Equals(appmode.DevelopmentMode)
		ctx = StoreLimitWhitelistDataInContext(ctx, canExceedLimit)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func StoreLimitWhitelistDataInContext(ctx context.Context, canExceedLimit bool) context.Context {
	return context.WithValue(ctx, limitWhitelistContextKey, canExceedLimit)
}

func CanExceedLimit(ctx context.Context) bool {
	whitelisted := ctx.Value(limitWhitelistContextKey)
	if whitelisted == nil {
		return false
	}
	return whitelisted.(bool)
}
