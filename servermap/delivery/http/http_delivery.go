package httpdelivery

import (
	"fmt"
	"github.com/Kichiyaki/appmode"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"

	"github.com/tribalwarshelp/map-generator/generator"

	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/servermap"
)

const (
	imageTTL = 2 * time.Hour / time.Second
	maxScale = 5
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
		return errors.New("cfg.MapUsecase is required")
	}
	h := &handler{cfg.MapUsecase, cfg.ServerUsecase}
	cfg.RouterGroup.GET("/map/:server", h.mapHandler)
	return nil
}

func (h *handler) mapHandler(c *gin.Context) {
	c.Header("Cache-Control", fmt.Sprintf(`public, max-age=%d`, imageTTL))

	srv, err := h.serverUsecase.GetByKey(c.Request.Context(), c.Param("server"))
	if err != nil {
		c.JSON(http.StatusNotFound, &gqlerror.Error{
			Message: err.Error(),
		})
		return
	}

	showBarbarian := c.Query("showBarbarian")
	largerMarkers := c.Query("largerMarkers")
	markersOnly := c.Query("markersOnly")
	markers, err := h.mapUsecase.GetMarkers(c.Request.Context(), servermap.GetMarkersConfig{
		Server:                  srv.Key,
		Tribes:                  c.Request.URL.Query()["tribe"],
		Players:                 c.Request.URL.Query()["player"],
		ShowBarbarianVillages:   showBarbarian == "true" || showBarbarian == "1",
		LargerMarkers:           largerMarkers == "true" || largerMarkers == "1",
		ShowOtherPlayerVillages: !(markersOnly == "true" || markersOnly == "1"),
		PlayerVillageColor:      c.Query("playerVillageColor"),
		BarbarianVillageColor:   c.Query("barbarianVillageColor"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, &gqlerror.Error{
			Message: err.Error(),
		})
		return
	}

	centerX, _ := strconv.Atoi(c.Query("centerX"))
	centerY, _ := strconv.Atoi(c.Query("centerY"))
	scale, _ := strconv.ParseFloat(c.Query("scale"), 32)
	if scale > maxScale {
		scale = maxScale
	}
	showGrid := c.Query("showGrid")
	showContinentNumbers := c.Query("showContinentNumbers")
	if err := generator.Generate(generator.Config{
		Markers:              markers,
		Destination:          c.Writer,
		ContinentGrid:        showGrid == "true" || showGrid == "1",
		ContinentNumbers:     showContinentNumbers == "true" || showContinentNumbers == "1",
		BackgroundColor:      c.Query("backgroundColor"),
		GridLineColor:        c.Query("gridLineColor"),
		ContinentNumberColor: c.Query("continentNumberColor"),
		MapSize:              srv.Config.Coord.MapSize,
		CenterX:              centerX,
		CenterY:              centerY,
		Scale:                float32(scale),
		Quality:              90,
		PNG:                  appmode.Equals(appmode.ProductionMode),
	}); err != nil {
		c.JSON(http.StatusBadRequest, &gqlerror.Error{
			Message: err.Error(),
		})
		return
	}
}
