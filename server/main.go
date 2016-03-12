package main

import (
  "./core"
  "fmt"
  "github.com/labstack/echo"
  // mw "github.com/labstack/echo/middleware"
)

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
  e.Run(":3030")
}
