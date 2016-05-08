package main

import (
	"challenger/server/core"
	"fmt"
	"github.com/labstack/echo"
	st "github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"os"
	"time"
)

const (
	httpAddr = "localhost:3000"
	tcpAddr  = "localhost:4000"
	udpAddr  = "localhost:5000"
	dbPath   = "./challenger.db"
)

func main() {
	// setup log system
	logfileName := "log/" + time.Now().Local().Format("2006-01-02-15-04-05") + ".log"
	f, err := os.OpenFile(logfileName, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println("error open log file", err)
		os.Exit(1)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	core.GetOptions()

	srv := core.NewSrv()
	go srv.Run(tcpAddr, udpAddr, dbPath)

	// setup echo
	ec := echo.New()
	ec.Static("/", "public")
	ec.Static("/api/asset/", "api_public")
	ec.Use(mw.Logger())
	ec.Get("/ws", st.WrapHandler(websocket.Handler(func(ws *websocket.Conn) {
		srv.ListenWebSocket(ws)
	})))
	ec.Post("/api/addteam", func(c echo.Context) error {
		return srv.AddTeam(c)
	})
	ec.Post("/api/resetqueue", func(c echo.Context) error {
		return srv.ResetQueue(c)
	})
	ec.Get("/api/history", func(c echo.Context) error {
		return srv.GetHistory(c)
	})
	ec.Post("/api/start_answer", func(c echo.Context) error {
		return srv.MatchStartAnswer(c)
	})
	ec.Post("/api/stop_answer", func(c echo.Context) error {
		return srv.MatchStopAnswer(c)
	})
	log.Println("listen http:", httpAddr)
	ec.Run(st.New(httpAddr))
}
