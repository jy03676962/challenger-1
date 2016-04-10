package core

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
)

var _ = log.Printf

func SetupRoute(e *echo.Echo) {
	e.Post("/login", echo.HandlerFunc(login))
}

func login(c echo.Context) error {
	return c.JSON(http.StatusOK, c.FormParams())
}
