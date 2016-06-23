package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var _ = log.Printf

type ReceiverInfo struct {
	ID    string `json:"id"`
	Idx   string `json:"idx"`
	Valid int    `json:"valid"`
}

type _laserMap map[string]*ReceiverInfo

type LaserPair struct {
	m _laserMap
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
	receiverKeys := make(map[string]string)
	for sender, info := range m {
		k := info.ID + ":" + info.Idx
		if s, ok := receiverKeys[k]; ok {
			log.Printf("got duplicated receiver:%v and %v to %v\n", s, sender, k)
			//os.Exit(1)
		} else {
			receiverKeys[k] = sender
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

func (l *LaserPair) GetValidReceivers() map[string]bool {
	ret := make(map[string]bool)
	for _, receiver := range l.m {
		if receiver.Valid == 1 {
			key := receiver.ID + ":" + receiver.Idx
			ret[key] = true
		}
	}
	return ret
}

func (l *LaserPair) GetValidSenders() map[string][]int {
	ret := make(map[string][]int)
	for k, receiver := range l.m {
		if receiver.Valid == 1 {
			li := strings.Split(k, ":")
			id := li[0]
			idx, _ := strconv.Atoi(li[1])
			value, ok := ret[id]
			if ok {
				ret[id] = append(value, idx)
			} else {
				ret[id] = []int{idx}
			}
		}
	}
	return ret
}

func (l *LaserPair) Record(key string, receiverID string, receiverIdx string, valid int) {
	info := ReceiverInfo{}
	info.ID = receiverID
	info.Idx = receiverIdx
	info.Valid = valid
	l.m[key] = &info
	l.Save()
}

func (l *LaserPair) RecordBrokens(brokens []string) {
	for _, broken := range brokens {
		li := strings.Split(broken, ":")
		for _, v := range l.m {
			if v.ID == li[0] && v.Idx == li[1] {
				v.Valid = 0
				break
			}
		}
	}
	l.Save()
}
