package core

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

var _ = log.Printf

func SetupRoute(e *echo.Echo) {
	e.Post("/api/addteam", echo.HandlerFunc(addTeam))
	e.Post("/api/resetqueue", echo.HandlerFunc(resetQueue))
}

func addTeam(c echo.Context) error {
	count, _ := strconv.Atoi(c.Param("count"))
	t, err := AddTeamToQueue(count)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func resetQueue(c echo.Context) error {
	err := ResetQueue()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
