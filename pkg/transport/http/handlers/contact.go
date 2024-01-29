package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zarldev/zarldotdev/view/contact"
)

type ContactHandler struct {
}

func NewContactHandler() *ContactHandler {
	return &ContactHandler{}
}

func (h *ContactHandler) RegisterRoutes(router *echo.Echo) {
	router.GET("/contact", h.Show)
}

func (h *ContactHandler) Show(c echo.Context) error {
	return render(c, contact.Show())
}
