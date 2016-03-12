package core

import (
  "golang.org/x/net/websocket"
  "log"
)

const (
  ROOM_CAPACITY = 4
)

type Server struct {
  clients   map[int]*Client
  addCh     chan *Client
  delCh     chan *Client
  sendAllCh chan map[string]interface{}
  doneCh    chan bool
  errCh     chan error
  room      *Room
  arg       *GameVar
}

func NewServer() *Server {
  clients := make(map[int]*Client)
  addCh := make(chan *Client)
  delCh := make(chan *Client)
  sendAllCh := make(chan map[string]interface{})
  doneCh := make(chan bool)
  errCh := make(chan error)
  arg := DefaultGameVar()

  return &Server{
    clients,
    addCh,
    delCh,
    sendAllCh,
    doneCh,
    errCh,
    nil,
    arg,
  }
}

func (s *Server) Add(c *Client) {
  s.addCh <- c
}

func (s *Server) Del(c *Client) {
  s.delCh <- c
}

func (s *Server) SendAll(msg map[string]interface{}) {
  s.sendAllCh <- msg
}

func (s *Server) Done() {
  s.doneCh <- true
}

func (s *Server) Err(err error) {
  s.errCh <- err
}

func (s *Server) handleMessage(msg map[string]interface{}, c *Client) {
  cmd := msg["cmd"].(string)
  switch cmd {
  case "login":
    data := make(map[string]interface{})
    data["cmd"] = "login"
    if s.room != nil {
      data["room"] = s.room
    }
    c.Write(data)
  case "createRoom":
    newRoom := Room{}
    newRoom.Hoster = msg["name"].(string)
    newRoom.MaxNum = ROOM_CAPACITY
    newRoom.Member = make([]string, 0)
    newRoom.Member = append(newRoom.Member, newRoom.Hoster)
    s.room = &newRoom
    data := make(map[string]interface{})
    data["cmd"] = "roomChanged"
    data["room"] = &newRoom
    data["arg"] = s.arg
    s.sendAll(data)
  case "joinRoom":
    if s.room != nil && len(s.room.Member) < s.room.MaxNum {
      s.room.Member = append(s.room.Member, msg["name"].(string))
      data := make(map[string]interface{})
      data["cmd"] = "roomChanged"
      data["room"] = s.room
      data["arg"] = s.arg
      s.sendAll(data)
    }
  case "startGame":
    data := make(map[string]interface{})
    data["cmd"] = "gameStarted"
    s.sendAll(data)
  default:
    c.Write(msg)
  }
}

func (s *Server) sendAll(msg map[string]interface{}) {
  for _, c := range s.clients {
    c.Write(msg)
  }
}

func (s *Server) OnConnected(ws *websocket.Conn) {
  defer func() {
    err := ws.Close()
    if err != nil {
      s.errCh <- err
    }
  }()

  client := NewClient(ws, s)
  s.Add(client)
  client.Listen()
}

func (s *Server) Start() {
  for {
    select {

    case c := <-s.addCh:
      log.Println("Added new client")
      s.clients[c.id] = c
      log.Println("Now", len(s.clients), "clients connected.")

    case c := <-s.delCh:
      log.Println("Delete client")
      delete(s.clients, c.id)

    case msg := <-s.sendAllCh:
      log.Println("Send all:", msg)
      s.sendAll(msg)

    case err := <-s.errCh:
      log.Println("Error:", err.Error())

    case <-s.doneCh:
      return
    }
  }
}
