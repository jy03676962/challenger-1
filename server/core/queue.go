package core

import (
	"container/list"
	"log"
	"strconv"
	"sync"
)

const initCursor = 2000

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
	for e, i := q.li.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		result[i] = e.Value.(*Team)
	}
	return result
}

//将拉去previousID之后count个team，previousID可以是0，表示从头开始拉取
//func GetTeamsFromQueue(previousID int, count int) ([]*Team, error) {
//if previousID < 0 {
//return nil, errors.New("previousID must >= 0")
//}
//q.lock.RLock()
//defer q.lock.RUnlock()
//if q.li.Len() == 0 {
//return nil, nil
//}
//var firstElement *list.Element = nil
//var remainTeamCount = q.li.Len()
//if previousID == 0 {
//firstElement = q.li.Front()
//} else {
//for e := q.li.Front(); e != nil; e = e.Next() {
//remainTeamCount -= 1
//team := e.Value.(*Team)
//if team.ID == previousID {
//firstElement = e.Next()
//break
//}
//}
//}
//if firstElement == nil {
//return nil, nil
//}
//resultCount := MinInt(remainTeamCount, count)
//result := make([]*Team, resultCount)
//for e, i := firstElement, 0; e != nil; e, i = e.Next(), i+1 {
//result[i] = e.Value.(*Team)
//}
//return result, nil
//}

func updateHallData() {
	msg := NewHubMap()
	msg.SetCmd("HallData")
	msg.Set("teams", GetAllTeamsFromQueue())
	socketInput := SocketInput{Broadcast: true, Group: SG_Admin, SocketMessage: msg}
	go func() {
		GetHub().SocketInputCh <- &socketInput
	}()
}
