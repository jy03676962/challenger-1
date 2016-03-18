package main

import (
	"challenger/server/core"
	"fmt"
	"github.com/labstack/echo"
	st "github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

//const HOST string = "172.16.10.59"
const HOST string = "localhost"

func main() {
	fmt.Println("start echo")
	srv := core.NewServer()
	go srv.Start()
	e := echo.New()
	e.Use(mw.Static("public"))
	e.Get("/ws", st.WrapHandler(websocket.Handler(func(ws *websocket.Conn) {
		srv.OnConnected(ws)
	})))
	e.Run(st.New(HOST + ":3030"))
}
