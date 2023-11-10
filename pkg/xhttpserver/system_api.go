package xhttpserver

import (
	"github.com/AltScore/gothic/v2/pkg/xopenapi"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/totemcaf/gollections/slices"
	"go.uber.org/zap"
)

type SystemModule struct {
	version string
	commit  string
}

func NewSystemApi(logger *zap.Logger, version string, commit string) *SystemModule {
	logger.Info("Starting system API")
	return &SystemModule{
		version: version,
		commit:  commit,
	}
}

func (s *SystemModule) RegisterHandlers(router xopenapi.EchoRouter) {
	router.GET("/health", s.health)
	router.GET("/routes", s.routes)
}

func (s *SystemModule) health(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"version": s.version,
		"commit":  s.commit,
		"status:": "OK",
	})
}

func (s *SystemModule) routes(c echo.Context) error {
	return c.JSON(http.StatusOK, slices.Map(c.Echo().Routes(), s.mapRoute))
}

type routeInfo struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (s *SystemModule) mapRoute(route *echo.Route) routeInfo {
	return routeInfo{
		Name:   cleanName(route.Name),
		Method: route.Method,
		Path:   route.Path,
	}
}

func cleanName(name string) string {
	segments := strings.Split(name, "/")
	if len(segments) == 1 || strings.Contains(segments[len(segments)-2], ".") {
		return segments[len(segments)-1]
	}

	return segments[len(segments)-2] + "." + segments[len(segments)-1]
}
