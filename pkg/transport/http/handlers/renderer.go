package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, comp templ.Component) error {
	return comp.Render(c.Request().Context(), c.Response())
}
