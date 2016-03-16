package main

import (
  "./core"
  "fmt"
  "github.com/labstack/echo"
  // mw "github.com/labstack/echo/middleware"
)

const HOST string = "172.16.10.59"

func main() {
  fmt.Println("start echo")
  srv := core.NewServer()
  go srv.Start()
  e := echo.New()
  e.Index("public/index.html")
  e.Static("/", "public")
  e.WebSocket("/ws", func(c *echo.Context) (err error) {
    srv.OnConnected(c.Socket())
    return
  })
  e.Run(HOST + ":3030")
}
