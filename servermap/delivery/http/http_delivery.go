package httpdelivery

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tribalwarshelp/map-generator/generator"

	"github.com/gin-gonic/gin"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/servermap"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	imageTTL = 2 * time.Hour / time.Second
)

type Config struct {
	RouterGroup   *gin.RouterGroup
	MapUsecase    servermap.Usecase
	ServerUsecase server.Usecase
}

type handler struct {
	mapUsecase    servermap.Usecase
	serverUsecase server.Usecase
}

func Attach(cfg Config) error {
	if cfg.MapUsecase == nil {
		return fmt.Errorf("cfg.MapUsecase cannot be nil")
	}
	h := &handler{cfg.MapUsecase, cfg.ServerUsecase}
	cfg.RouterGroup.GET("/map/:server", h.mapHandler)
	return nil
}

func (h *handler) mapHandler(c *gin.Context) {
	c.Header("Cache-Control", fmt.Sprintf(`public, must-revalidate, max-age=%d, s-maxage=%d`, imageTTL, imageTTL))

	server, err := h.serverUsecase.GetByKey(c.Request.Context(), c.Param("server"))
	if err != nil {
		c.JSON(http.StatusNotFound, &gqlerror.Error{
			Message: err.Error(),
		})
		return
	}

	markers, err := h.mapUsecase.GetMarkers(c.Request.Context(), servermap.GetMarkersConfig{
		Server:                  server.Key,
		Tribes:                  c.Request.URL.Query()["tribe"],
		Players:                 c.Request.URL.Query()["player"],
		ShowBarbarianVillages:   c.Query("showBarbarian") == "true",
		ShowOtherPlayerVillages: !(c.Query("onlyMarkers") == "true"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, &gqlerror.Error{
			Message: err.Error(),
		})
		return
	}

	if err := generator.Generate(generator.Config{
		Markers:              markers,
		Destination:          c.Writer,
		ContinentGrid:        c.Query("showGrid") == "true",
		ContinentNumbers:     c.Query("showContinentNumbers") == "true",
		BackgroundColor:      c.Query("backgroundColor"),
		GridLineColor:        c.Query("gridLineColor"),
		ContinentNumberColor: c.Query("continentNumberColor"),
		MapSize:              server.Config.Coord.MapSize,
	}); err != nil {
		c.JSON(http.StatusBadRequest, &gqlerror.Error{
			Message: err.Error(),
		})
		return
	}
}
