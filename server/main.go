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

const API string = "localhost:4040"
const HOST string = "localhost:3030"

func main() {
	logfileName := time.Now().Local().Format("2006-01-02-15-04-05") + ".log"
	f, err := os.OpenFile(logfileName, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println("error open log file", err)
		os.Exit(1)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	hub := core.NewHub()
	log.Println("start listen websocket:", HOST)
	srv := core.NewServer(hub)
	api := core.NewTCPServer(API, hub)
	match := core.NewMatch(hub)
	go match.Run()
	go srv.Run()
	go api.Run()
	e := echo.New()
	e.Static("/", "public")
	e.Use(mw.Logger())
	e.Get("/ws", st.WrapHandler(websocket.Handler(func(ws *websocket.Conn) {
		srv.OnConnected(ws)
	})))
	e.Get("/api", st.WrapHandler(websocket.Handler(func(ws *websocket.Conn) {
		srv.OnApiConnected(ws)
	})))
	core.SetupRoute(e)
	e.Run(st.New(HOST))
}
