package main

import (
	"challenger/server/core"
	"github.com/labstack/echo"
	st "github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
	"log"
)

const API string = "localhost:4040"
const HOST string = "localhost:3030"

func main() {
	hub := core.NewHub()
	log.Println("start listen websocket:", HOST)
	srv := core.NewServer(hub)
	api := core.NewTCPServer(API, hub)
	match := core.NewMatch(hub)
	go match.Run()
	go srv.Run()
	go api.Run()
	e := echo.New()
	e.Use(mw.Static("public"))
	e.Get("/ws", st.WrapHandler(websocket.Handler(func(ws *websocket.Conn) {
		srv.OnConnected(ws)
	})))
	e.Get("/api", st.WrapHandler(websocket.Handler(func(ws *websocket.Conn) {
		srv.OnApiConnected(ws)
	})))
	e.Run(st.New(HOST))
}
