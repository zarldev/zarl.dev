package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zarldev/zarldotdev/pkg/repo"
	"github.com/zarldev/zarldotdev/view/component/claps"
)

type ClapsHandler struct {
	Repo *repo.ClapsRepository
}

func NewClapsHandler(config repo.Config) (*ClapsHandler, error) {
	conn, err := repo.NewConnection(config)
	if err != nil {
		return nil, err
	}
	cr, err := repo.NewClapsRepository(conn)
	if err != nil {
		return nil, err
	}
	return &ClapsHandler{
		Repo: cr,
	}, nil
}

func (h *ClapsHandler) RegisterRoutes(router *echo.Echo) {
	clapsGroup := router.Group("/claps")
	clapsGroup.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
	clapsGroup.GET("/:id", h.GetInitialClaps)
	clapsGroup.POST("/:id", h.AddClap)
}

func (h *ClapsHandler) GetInitialClaps(c echo.Context) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	count, err := h.Repo.Get(id)
	if err != nil {
		return err
	}

	cookies := c.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == clapID(id) {
			break
		}
	}
	return render(c, claps.Clap(id, count))
}

func idParam(c echo.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *ClapsHandler) AddClap(c echo.Context) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	count, err := h.Repo.Increment(id)
	if err != nil {
		return err
	}
	c.SetCookie(&http.Cookie{
		Name:     clapID(id),
		Value:    "true",
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		Secure:   false,
		HttpOnly: false,
	})
	return render(c, claps.ClappedWithHeader(count))
}

func clapID(id int) string {
	return fmt.Sprintf("clapped_%d", id)
}
