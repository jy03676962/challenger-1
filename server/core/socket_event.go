package core

type SocketEvent struct {
  SocketMessage map[string]interface{}
  Client        *Client
}
