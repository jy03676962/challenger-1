package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var _ = log.Printf

type ReceiverInfo struct {
	ID    string `json:"id"`
	Idx   string `json:"idx"`
	Valid int    `json:"valid"`
}

type _laserMap map[string]ReceiverInfo

type LaserPair struct {
	m       _laserMap
	brokens map[string]int
}

var laserPair = loadLaserPair()

func GetLaserPair() *LaserPair {
	return laserPair
}

func loadLaserPair() *LaserPair {
	m := make(_laserMap)
	b, e := ioutil.ReadFile("./laser.json")
	if os.IsNotExist(e) {
		return newLaserPair(m)
	}
	if e != nil {
		log.Printf("parse laser pair error:%v\n", e.Error())
		os.Exit(1)
	}
	e = json.Unmarshal(b, &m)
	if e != nil {
		log.Printf("parse laser pair error:%v\n", e.Error())
		os.Exit(1)
	}
	return newLaserPair(m)
}

func newLaserPair(m _laserMap) *LaserPair {
	lp := LaserPair{}
	lp.m = m
	lp.brokens = make(map[string]int)
	for _, rcvrs := range m {
		if rcvrs.Valid == 0 {
			key := rcvrs.ID + ":" + rcvrs.Idx
			lp.brokens[key] = 1
		}
	}
	return &lp
}

func (l *LaserPair) Save() {
	b, _ := json.Marshal(l.m)
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	ioutil.WriteFile("./laser.json", out.Bytes(), 0640)
}

func (l *LaserPair) Record(key string, receiverID string, receiverIdx string) {
	info := ReceiverInfo{}
	info.ID = receiverID
	info.Idx = receiverIdx
	info.Valid = 1
	l.m[key] = info
	l.Save()
}
