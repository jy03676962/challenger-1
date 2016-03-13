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

type ArenaWall struct {
  P1 P
  P2 P
}

type W ArenaWall
