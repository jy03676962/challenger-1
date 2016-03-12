package core

// the skeleton of this file is borrowed from https://github.com/golang-samples/websocket

import (
  "golang.org/x/net/websocket"
  "log"
)

type Server struct {
  clients   map[int]*Client
  addCh     chan *Client
  delCh     chan *Client
  sendAllCh chan map[string]interface{}
  doneCh    chan bool
  errCh     chan error
  match     *Match
}

func NewServer() *Server {
  clients := make(map[int]*Client)
  addCh := make(chan *Client)
  delCh := make(chan *Client)
  sendAllCh := make(chan map[string]interface{})
  doneCh := make(chan bool)
  errCh := make(chan error)

  return &Server{
    clients,
    addCh,
    delCh,
    sendAllCh,
    doneCh,
    errCh,
    nil,
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
  name := ""
  if msg["name"] != nil {
    name = msg["name"].(string)
  }
  switch cmd {
  case "login":
    data := make(map[string]interface{})
    data["cmd"] = "login"
    c.SetUsername(name)
    // if we have a match then tell client
    if s.match != nil {
      data["match"] = s.match
    }
    c.Write(data)
  case "createMatch":
    newMatch := NewMatch()
    newMatch.Hoster = name
    newMatch.AddMember(name)
    s.match = newMatch
    data := make(map[string]interface{})
    data["cmd"] = "matchChanged"
    data["match"] = newMatch
    s.sendAll(data)
  case "joinMatch":
    if s.match != nil {
      if s.match.AddMember(name) {
        data := make(map[string]interface{})
        data["cmd"] = "matchChanged"
        data["match"] = s.match
        s.sendAll(data)
      }
    }
  case "startMatch":
    s.match.Start()
    data := make(map[string]interface{})
    data["cmd"] = "matchChanged"
    data["match"] = s.match
    s.sendAll(data)
  default:
    c.Write(msg)
  }
}

func (s *Server) getClient(name string) *Client {
  for _, c := range s.clients {
    if c.GetUsername() == name {
      return c
    }
  }
  return nil
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
