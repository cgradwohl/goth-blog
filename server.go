package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
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

	errorPage := fmt.Sprintf("public/error/%d.html", code)

	fileError := c.File(errorPage)
	if fileError != nil {
		c.Logger().Error(fileError)
	}
}

type Config struct {
	addr string
	tls  *tls.Config
}

func NewConfig() Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	if env == "production" {
		autoTLSManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/var/www/.cache"),
			HostPolicy: autocert.HostWhitelist("<DOMAIN>"),
		}

		return Config{
			addr: ":443",
			tls: &tls.Config{
				GetCertificate: autoTLSManager.GetCertificate,
				NextProtos:     []string{acme.ALPNProto},
			},
		}
	}

	return Config{
		addr: ":3000",
	}
}

func NewRouter() *echo.Echo {
	router := echo.New()
	router.HTTPErrorHandler = customHTTPErrorHandler

	router.Use(middleware.Recover())
	router.Use(middleware.Logger())
	router.Static("/", "public")

	router.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, page.Home())
	})

	router.GET("/foo", func(c echo.Context) error {
		return Render(c, http.StatusOK, page.Foo())
	})

	return router
}

type Server struct {
	addr    string
	handler *echo.Echo
	tls     *tls.Config
}

func NewServer() (*Server, error) {
	config := NewConfig()
	router := NewRouter()

	return &Server{
		addr:    config.addr,
		handler: router,
		tls:     config.tls,
	}, nil
}

func (server *Server) start() error {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	httpServer := &http.Server{
		Addr:    server.addr,
		Handler: server.handler,
	}

	if env == "production" {
		autoTLSManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/var/www/.cache"),
			HostPolicy: autocert.HostWhitelist("<DOMAIN>"),
		}

		server.tls = &tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		}

		httpServer.TLSConfig = server.tls
		// TODO: add certificates for production
		log.Printf("hello creature...listening on %s in Prod", server.addr)
		return httpServer.ListenAndServeTLS("", "")
	}

	log.Printf("hello creature...listening on %s in Dev", server.addr)
	return httpServer.ListenAndServe()
}

func main() {
	server, err := NewServer()
	if err != nil {
		panic(err)
	}

	err = server.start()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
