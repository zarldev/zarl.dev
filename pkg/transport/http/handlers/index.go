package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zarldev/zarldotdev/view/index"
)

type IndexHandler struct {
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.Index)
}

func (h *IndexHandler) Index(c echo.Context) error {
	return render(c, index.Welcome())
}
