package core

type ArenaPosition struct {
	X int
	Y int
}

type P ArenaPosition

type RealPosition struct {
	X float64
	Y float64
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

type Button struct {
	Id string `json:"id"`
	R  Rect   `json:"r"`
}
