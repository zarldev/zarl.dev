package handlers

import "github.com/labstack/echo/v4"

type StaticHandler struct {
}

func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

func (h *StaticHandler) RegisterRoutes(router *echo.Echo) {
	router.Static("/static", "assets")
}
