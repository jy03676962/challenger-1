package core

import (
	"log"
	"math"
)

const laserSize = 10

var _ = log.Printf

var catchByPos = false

type LaserLine struct {
	ID    string
	Index int
	P     int
}

type Laser struct {
	IsPause              bool `json:"isPause"`
	DisplayP             RP   `json:"displayP"`
	DisplayP2            RP   `json:"displayP2"`
	player               *Player
	dest                 int
	pathMap              map[int]int
	p                    int
	p2                   int
	match                *Match
	pauseTime            float64
	elaspedSinceLastMove float64
	lines                []LaserLine
	startupLines         []LaserLine
	startupingIndex      int
	closed               bool
}

func NewLaser(p P, player *Player, match *Match) *Laser {
	l := Laser{}
	l.IsPause = true
	l.player = player
	l.dest = -1
	l.match = match
	l.pathMap = make(map[int]int)
	l.p = GetOptions().TilePosToInt(p)
	l.p2 = -1
	l.convertDisplay()
	l.elaspedSinceLastMove = GetOptions().LaserSpeed
	l.lines = make([]LaserLine, 0)
	l.startupLines = make([]LaserLine, 0)
	infos := GetOptions().mainArduinoInfosByPos(l.p)
	for _, info := range infos {
		for i := 0; i < info.LaserNum; i++ {
			l.startupLines = append(l.startupLines, LaserLine{info.ID, i, l.p})
		}
	}
	l.startupingIndex = 0
	l.closed = true
	return &l
}

func (l *Laser) Pause(t float64) {
	if l.isStartuping() {
		return
	}
	l.IsPause = true
	l.pauseTime = math.Max(t, l.pauseTime)
	l.doClose()
}

func (l *Laser) IsFollow(cid string) bool {
	return l.player.ControllerID == cid
}

func (l *Laser) Close() {
	l.closed = true
	l.doClose()
}

func (l *Laser) IsTouched(changes *[]laserInfoChange) (touched bool, p int) {
	p = 0
	touched = false
	if l.IsPause || l.isStartuping() {
		return
	}
	for _, change := range *changes {
		for _, line := range l.lines {
			if line.ID == change.id && line.Index == change.idx {
				p = line.P
				touched = true
				return
			}
		}
	}
	return
}

func (l *Laser) Tick(dt float64) {
	if l.closed {
		return
	}
	if l.IsPause {
		l.pauseTime -= dt
		if l.pauseTime <= 0 {
			l.IsPause = false
			l.pauseTime = 0
			l.doOpen()
		}
		return
	}
	opt := GetOptions()
	l.elaspedSinceLastMove += dt
	interval := opt.laserMoveInterval(l.match.Energy)
	if l.elaspedSinceLastMove < interval {
		return
	}
	l.elaspedSinceLastMove = 0
	if l.isStartuping() {
		line := l.startupLines[l.startupingIndex]
		l.match.openLaser(line.ID, line.Index)
		l.lines = append(l.lines, line)
		l.startupingIndex += 1
	} else {
		next := l.findPath()
		if l.p2 < 0 && l.p == next {
			return
		}
		replaceIdx := -1
		notInNext := 0
		for i, line := range l.lines {
			if line.P != next {
				notInNext += 1
				if replaceIdx < 0 {
					replaceIdx = i
				}
			}
		}
		if replaceIdx >= 0 {
			infos := opt.mainArduinoInfosByPos(next)
			for _, info := range infos {
				for i := 0; i < info.LaserNum; i++ {
					if !l.contains(info.ID, i) {
						line := l.lines[replaceIdx]
						l.match.closeLaser(line.ID, line.Index)
						l.match.openLaser(info.ID, i)
						l.lines[replaceIdx] = LaserLine{info.ID, i, next}
						if notInNext == 1 {
							l.p = next
							l.p2 = -1
						} else {
							l.p2 = next
						}
						l.convertDisplay()
						return
					}
				}
			}
		}
	}
}

func (l *Laser) doClose() {
	for _, line := range l.lines {
		l.match.closeLaser(line.ID, line.Index)
	}
}

func (l *Laser) doOpen() {
	for _, line := range l.lines {
		l.match.openLaser(line.ID, line.Index)
	}
}

func (l *Laser) convertDisplay() {
	opt := GetOptions()
	y := l.p / opt.ArenaWidth
	x := l.p % opt.ArenaWidth
	y = opt.ArenaHeight - 1 - y
	l.DisplayP = opt.RealPosition(P{x, y})
	if l.p2 >= 0 {
		y := l.p2 / opt.ArenaWidth
		x := l.p2 % opt.ArenaWidth
		y = opt.ArenaHeight - 1 - y
		l.DisplayP2 = opt.RealPosition(P{x, y})
	} else {
		l.DisplayP2 = RP{-1, -1}
	}
}

func (l *Laser) contains(id string, idx int) bool {
	for _, line := range l.lines {
		if line.ID == id && line.Index == idx {
			return true
		}
	}
	return false
}

func (l *Laser) isStartuping() bool {
	return l.startupingIndex < len(l.startupLines)
}

func (l *Laser) findPath() int {
	opt := GetOptions()
	l.fillPath()
	pp1 := opt.Conv(l.p)
	if l.p2 >= 0 {
		pp2 := opt.Conv(l.p2)
		if l.pathMap[pp2] <= l.pathMap[pp1] {
			return l.p2
		} else {
			return l.p
		}
	}
	next, min := pp1, l.pathMap[pp1]
	for _, i := range opt.TileAdjacency[pp1] {
		if l.pathMap[i] < min {
			min = l.pathMap[i]
			next = i
		}
	}
	return opt.Conv(next)
}

func (l *Laser) fillPath() {
	opt := GetOptions()
	dest := opt.TilePosToInt(l.player.tilePos)
	if l.dest == dest {
		return
	}
	l.dest = dest
	for i := 0; i < opt.ArenaWidth*opt.ArenaHeight; i++ {
		l.pathMap[i] = 10000
	}
	var fill func(x int, v int)
	fill = func(x int, v int) {
		l.pathMap[x] = v
		for _, i := range opt.TileAdjacency[x] {
			if l.pathMap[i] > v+1 {
				fill(i, v+1)
			}
		}
	}
	fill(opt.Conv(l.dest), 0)
}
