package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"wohlburger.io/layout"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func main() {
	e := echo.New()
	fmt.Println("hello creature ...")

	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, layout.MainLayout("Chris"))
	})

	e.Start(":8080")
}
