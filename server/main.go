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
	host        = "localhost"
	httpAddr    = host + ":3000"
	tcpAddr     = host + ":4000"
	udpAddr     = host + ":5000"
	dbPath      = "./challenger.db"
	isSimulator = true
)

func main() {
	// setup log system
	log.Println("start server")
	logfileName := "log/" + time.Now().Local().Format("2006-01-02-15-04-05") + ".log"
	os.Mkdir("log", 0777)
	f, err := os.OpenFile(logfileName, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println("error open log file", err)
		os.Exit(1)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	log.Println("setup log system done")

	core.GetOptions()
	core.GetSurvey()

	log.Println("reading cfg done")

	srv := core.NewSrv(isSimulator)
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
	ec.Get("/api/survey", func(c echo.Context) error {
		return srv.GetSurvey(c)
	})
	ec.Post("/api/answer", func(c echo.Context) error {
		return srv.UpdateQuestionInfo(c)
	})
	ec.Get("/api/answering", func(c echo.Context) error {
		return srv.GetAnsweringMatchData(c)
	})
	ec.Post("/api/update_player", func(c echo.Context) error {
		return srv.UpdatePlayerData(c)
	})
	ec.Post("/api/update_match", func(c echo.Context) error {
		return srv.UpdateMatchData(c)
	})
	ec.Get("/api/sender_list", func(c echo.Context) error {
		return srv.GetMainArduinoList(c)
	})
	log.Println("listen http:", httpAddr)
	ec.Run(st.New(httpAddr))
}
