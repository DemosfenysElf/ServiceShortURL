package router

import (
	"net"
	"net/http"

	"github.com/labstack/echo"
)

func (s serverShortener) MWCheakerIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.Cfg.TrustedSubnet == "" {
			c.Response().WriteHeader(http.StatusForbidden)
			return nil
		}

		getIP := c.Request().Header.Get("X-Real-IP")
		if getIP == "" {
			c.Response().WriteHeader(http.StatusForbidden)
			return nil
		}

		_, cfgCIDR, err := net.ParseCIDR(s.Cfg.TrustedSubnet)
		if err != nil {
			c.Response().WriteHeader(http.StatusForbidden)
			return nil
		}

		netIP := net.ParseIP(getIP)
		if !cfgCIDR.Contains(netIP) {
			c.Response().WriteHeader(http.StatusForbidden)
			return nil
		}

		return next(c)
	}
}
