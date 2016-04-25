package core

import (
	"container/list"
	"log"
	"strconv"
	"sync"
)

const (
	initCursor     = 2000
	singleWaitTime = 360
)

type TeamStatus int

const (
	TS_Waiting  TeamStatus = iota
	TS_Prepare  TeamStatus = iota
	TS_Playing  TeamStatus = iota
	TS_After    TeamStatus = iota
	TS_Finished TeamStatus = iota
)

var _ = log.Printf

var q = newQueue()

type Team struct {
	Size       int        `json:"size"`
	ID         string     `json:"id"`
	DelayCount int        `json:"delayCount"`
	Status     TeamStatus `json:"status"`
	WaitTime   int        `json:"waitTime"`
}

type Queue struct {
	li   *list.List
	dict map[string]*list.Element
	cur  int
	lock *sync.RWMutex
}

func newQueue() *Queue {
	q := Queue{}
	q.li = list.New()
	q.dict = make(map[string]*list.Element)
	q.cur = initCursor
	q.lock = new(sync.RWMutex)
	return &q
}

func AddTeamToQueue(teamSize int) (*Team, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	defer updateHallData()
	q.cur += 1
	id := strconv.Itoa(q.cur)
	t := Team{Size: teamSize, ID: id, Status: TS_Waiting}
	element := q.li.PushBack(&t)
	q.dict[id] = element
	return &t, nil
}

func ResetQueue() error {
	q.lock.Lock()
	defer q.lock.Unlock()
	defer updateHallData()
	q.li.Init()
	q.dict = make(map[string]*list.Element)
	q.cur = initCursor
	return nil
}

func EnterPrepare() error {
	return nil
}

func EnterPlay() error {
	return nil
}

func TeamCutLine(teamID string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	defer updateHallData()
	element := q.dict[teamID]
	team := element.Value.(*Team)
	if team.Status != TS_Waiting {
		return
	}
	for e := q.li.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Team)
		if t.Status == TS_Waiting {
			log.Printf("becut:%v\n", t)
			q.li.MoveBefore(element, e)
			return
		}
	}
}

func TeamRemove(teamID string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	defer updateHallData()
	element := q.dict[teamID]
	team := element.Value.(*Team)
	if team.Status != TS_Waiting {
		return
	}
	delete(q.dict, teamID)
	q.li.Remove(element)
}

func GetAllTeamsFromQueueWithLock() []*Team {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return GetAllTeamsFromQueue()
}

func GetAllTeamsFromQueue() []*Team {
	result := make([]*Team, q.li.Len())
	waitTime := 0
	for e, i := q.li.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		team := e.Value.(*Team)
		if team.Status == TS_Waiting {
			waitTime += singleWaitTime
			team.WaitTime = waitTime
		}
		result[i] = team
	}
	return result
}

func updateHallData() {
	msg := NewHubMap()
	msg.SetCmd("HallData")
	msg.Set("teams", GetAllTeamsFromQueue())
	socketInput := SocketInput{Broadcast: true, Group: SG_Admin, SocketMessage: msg}
	go func() {
		GetHub().SocketInputCh <- &socketInput
	}()
}
