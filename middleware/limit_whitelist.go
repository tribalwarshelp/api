package middleware

import (
	"context"
	"github.com/Kichiyaki/appmode"
	"net"

	"github.com/gin-gonic/gin"
)

var limitWhitelistContextKey ContextKey = "limitWhitelist"

type NetworksAndIps struct {
	Networks []*net.IPNet
	Ips      []net.IP
}

func (networksAndIps NetworksAndIps) Contains(ip net.IP) bool {
	for _, whitelistedIP := range networksAndIps.Ips {
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

func (cfg LimitWhitelistConfig) GetNetworksAndIps() NetworksAndIps {
	networks := []*net.IPNet{}
	ips := []net.IP{}
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
	return NetworksAndIps{
		Networks: networks,
		Ips:      ips,
	}
}

func LimitWhitelist(cfg LimitWhitelistConfig) gin.HandlerFunc {
	networksAndIps := cfg.GetNetworksAndIps()
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
