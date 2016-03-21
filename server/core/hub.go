package core

type HubCmd interface {
	GetCmd() string
	SetCmd(cmd string)
}

type HubMap struct {
	m map[string]interface{}
}

func NewHubMap() *HubMap {
	ret := HubMap{}
	ret.m = make(map[string]interface{})
	return &ret
}

func NewHubMapWithData(m map[string]interface{}) *HubMap {
	ret := HubMap{}
	ret.m = m
	return &ret
}

func (hm *HubMap) Data() map[string]interface{} {
	return hm.m
}

func (hm *HubMap) Get(key string) interface{} {
	if v, ok := hm.m[key]; ok {
		return v
	}
	return nil
}

func (hm *HubMap) GetStr(key string) string {
	if v, ok := hm.m[key]; ok {
		return v.(string)
	}
	return ""
}

func (hm *HubMap) Set(key string, value interface{}) {
	hm.m[key] = value
}

func (hm *HubMap) GetCmd() string {
	return hm.GetStr("cmd")
}

func (hm *HubMap) SetCmd(v string) {
	hm.Set("cmd", v)
}

type Hub struct {
	MatchInputCh   chan *HubMap
	MatchOutputCh  chan *HubMap
	ServerQuitCh   chan struct{}
	SocketOutputCh chan *SocketOutput
	SocketInputCh  chan *SocketInput
	Options        *MatchOptions
}

func NewHub() *Hub {
	hub := Hub{}
	hub.MatchInputCh = make(chan *HubMap)
	hub.MatchOutputCh = make(chan *HubMap)
	hub.ServerQuitCh = make(chan struct{})
	hub.SocketOutputCh = make(chan *SocketOutput)
	hub.SocketInputCh = make(chan *SocketInput)
	hub.Options = DefaultMatchOptions()
	return &hub
}
