package core

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
)

var _ = log.Printf

func SetupRoute(e *echo.Echo) {
	e.Post("/login", echo.HandlerFunc(login))
	e.Post("/latest", echo.HandlerFunc(latest))
}

func login(c echo.Context) error {
	return c.JSON(http.StatusOK, c.FormParams())
}

func latest(c echo.Context) error {
	return nil
}
