package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zarldev/zarldotdev/view/about"
)

type AboutHandler struct {
}

func NewAboutHandler() *AboutHandler {
	return &AboutHandler{}
}

func (h *AboutHandler) RegisterRoutes(router *echo.Echo) {
	router.GET("/about", h.AboutShow)
}

func (h *AboutHandler) AboutShow(c echo.Context) error {
	return render(c, about.Show())
}
