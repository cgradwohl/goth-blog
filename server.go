package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"wohlburger.io/page"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError

	httpError, ok := err.(*echo.HTTPError)
	if ok {
		code = httpError.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("pubic/error/%d.html", code)

	fileError := c.File(errorPage)
	if fileError != nil {
		c.Logger().Error(err)
	}
}

func main() {
	e := echo.New()
	fmt.Println("hello creature ...")

	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, page.Home())
	})

	e.Start(":8080")
}
