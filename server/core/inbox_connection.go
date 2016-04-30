package core

import (
	"bufio"
	"encoding/json"
	"golang.org/x/net/websocket"
	"net"
	"strings"
	"sync"
)

type InboxConnection interface {
	ReadJSON(v *InboxMessage) error
	WriteJSON(v *InboxMessage) error
	Close() error
	Accept(addr InboxAddress) bool
}

type InboxTcpConnection struct {
	conn *net.TCPConn
	r    *bufio.Reader
	id   string
}

func NewInboxTcpConnection(conn *net.TCPConn) *InboxTcpConnection {
	tcp := InboxTcpConnection{conn: conn}
	tcp.r = bufio.NewReader(conn)
	return &tcp
}

func (tcp *InboxTcpConnection) Close() error {
	return tcp.conn.Close()
}

func (tcp *InboxTcpConnection) ReadJSON(v *InboxMessage) error {
	b, e := tcp.r.ReadBytes(60) // tcp message frame start with '<'
	if e != nil {
		return e
	}
	b, e = tcp.r.ReadBytes(62) // tcp message frame end with '>'
	if e != nil {
		return e
	}
	if len(b) == 1 { // only has '>' delimiter
		return nil
	}
	if b[0] == 123 { // first byte is '{', json encoding frame
		e := json.Unmarshal(b[:len(b)-1], &v.Data)
		if e == nil {
			if len(tcp.id) > 0 {
				v.Set("ID", tcp.id)
				v.Set("TYPE", InboxAddressTypeArduinoDevice)
			}
		}
		return e
	} else { // parse heart beat frame
		parseTcpHB(string(b), v)
		v.SetCmd("hb")
		v.Set("TYPE", InboxAddressTypeArduinoDevice)
		if id := v.GetStr("ID"); id != "" {
			tcp.id = id
		}
		return nil
	}
}

// Tcp HB format is [key1]value1[key2]value2
func parseTcpHB(hb string, v *InboxMessage) {
	kvs := strings.Split(hb, "[")
	for _, s := range kvs {
		kv := strings.Split(s, "]")
		if len(kv) == 2 {
			v.Set(kv[0], kv[1])
		}
	}
}

func (tcp *InboxTcpConnection) WriteJSON(v *InboxMessage) error {
	b, e := v.Marshal()
	if e != nil {
		return e
	}

	buf := make([]byte, len(b)+2)
	for i := 1; i < len(buf)-1; i++ {
		buf[i] = b[i-1]
	}
	buf[0] = 60
	buf[len(buf)-1] = 62
	_, e = tcp.conn.Write(buf)
	return e
}

func (tcp *InboxTcpConnection) Accept(addr InboxAddress) bool {
	if addr.Type != InboxAddressTypeAdminDevice {
		return false
	}
	return addr.ID == "" || addr.ID == tcp.id
}

type InboxUdpConnection struct {
	conn *net.UDPConn
	dict map[string]*net.UDPAddr
	lock *sync.RWMutex
}

func NewInboxUdpConnection(conn *net.UDPConn) *InboxUdpConnection {
	u := InboxUdpConnection{conn: conn}
	u.dict = make(map[string]*net.UDPAddr)
	u.lock = new(sync.RWMutex)
	return &u
}

func (udp *InboxUdpConnection) Close() error {
	return udp.conn.Close()
}

func (udp *InboxUdpConnection) ReadJSON(v *InboxMessage) error {
	buf := make([]byte, 1024)
	n, addr, err := udp.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	cmdLen := 11
	if n >= cmdLen {
		d := buf[:cmdLen]
		id := string(d[3:6])
		v.Set("head", string(d[:3]))
		v.Set("ID", id)
		v.Set("loc", string(d[6:9]))
		v.Set("status", string(d[9:]))
		v.Set("TYPE", InboxAddressTypeWearableDevice)
		udp.lock.Lock()
		udp.dict[id] = addr
		udp.lock.Unlock()
	}
	return nil
}
func (udp *InboxUdpConnection) WriteJSON(v *InboxMessage) error {
	str := v.GetStr("head") + v.GetStr("id") + v.GetStr("cmd")
	udp.lock.RLock()
	addr := udp.dict[v.GetStr("id")]
	udp.lock.RUnlock()
	if addr != nil {
		_, e := udp.conn.WriteToUDP([]byte(str), addr)
		return e
	}
	return nil
}

func (udp *InboxUdpConnection) Accept(addr InboxAddress) bool {
	if addr.Type != InboxAddressTypeWearableDevice {
		return false
	}
	if addr.ID == "" {
		return true
	}
	udp.lock.RLock()
	defer udp.lock.RUnlock()
	_, ok := udp.dict[addr.ID]
	return ok
}

type InboxWsConnection struct {
	conn *websocket.Conn
	t    InboxAddressType
	id   string
	l    *sync.RWMutex
}

func NewInboxWsConnection(conn *websocket.Conn) *InboxWsConnection {
	return &InboxWsConnection{conn: conn, l: new(sync.RWMutex)}
}

func (ws *InboxWsConnection) Close() error {
	return ws.conn.Close()
}

func (ws *InboxWsConnection) ReadJSON(v *InboxMessage) error {
	e := websocket.JSON.Receive(ws.conn, &v.Data)
	if e != nil {
		return e
	}
	if v.GetCmd() == "init" {
		t := v.GetStr("TYPE")
		var tt InboxAddressType
		id := ""
		if t == "admin" {
			tt = InboxAddressTypeAdminDevice
		} else if t == "postgame" {
			tt = InboxAddressTypePostgameDevice
		} else if t == "simulator" {
			tt = InboxAddressTypePostgameDevice
		}
		ws.setAddressInfo(v.GetStr("ID"), tt)
	}
	if id, t, b := ws.getAddressInfo(); b {
		v.Set("ID", id)
		v.Set("TYPE", t)
	}
	return nil
}

func (ws *InboxWsConnection) WriteJSON(v *InboxMessage) error {
	return websocket.JSON.Send(ws.conn, v.Data)
}

func (ws *InboxWsConnection) Accept(addr InboxAddress) bool {
	if id, t, b := ws.getAddressInfo(); b {
		if t != addr.Type {
			return false
		}
		return addr.ID == "" || id == addr.ID
	}
	return false
}

func (ws *InboxWsConnection) getAddressInfo() (id string, t InboxAddressType, hasInfo bool) {
	id, t, hasInfo = "", InboxAddressTypeUnknown, false
	ws.l.RLock()
	defer ws.l.RUnlock()
	if len(ws.id) > 0 {
		id, t, hasInfo = ws.id, ws.t, true
	}
	return
}

func (ws *InboxWsConnection) setAddressInfo(id string, t InboxAddressType) {
	ws.l.Lock()
	defer ws.l.Unlock()
	ws.id, ws.t = id, t
}
