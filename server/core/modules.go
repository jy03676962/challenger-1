package core

type ArenaPosition struct {
  X int
  Y int
}

type P ArenaPosition

type RealPosition struct {
  X float32
  Y float32
}

type RP RealPosition

type Rect struct {
  X float64
  Y float64
  W float64
  H float64
}

type ArenaWall struct {
  P1 P
  P2 P
}

type W ArenaWall

type SocketEvent struct {
  SocketMessage map[string]interface{}
  Client        *Client
}
