package core

import (
	"log"
	"math"
)

var _ = log.Printf

type Laser struct {
	Pos     RP   `json:"pos"`
	IsPause bool `json:"isPause"`
	//private
	player    *Player
	next      int
	dest      int
	pathMap   map[int]int
	p         int
	match     *Match
	pauseTime float64
}

func NewLaser(p P, player *Player, match *Match) *Laser {
	l := Laser{}
	l.IsPause = true
	l.player = player
	l.next = -1
	l.dest = -1
	l.match = match
	l.pathMap = make(map[int]int)
	l.p = l.getOpt().TilePosToInt(p)
	l.Pos = l.getOpt().RealPosition(p)
	return &l
}

func (l *Laser) Pause(t float64) {
	l.IsPause = true
	l.pauseTime = math.Max(t, l.pauseTime)
}

func (l *Laser) IsFollow(cid string) bool {
	return l.player.ControllerID == cid
}

func (l *Laser) Tick(dt float64) {
	if l.IsPause {
		l.pauseTime -= dt
		if l.pauseTime <= 0 {
			l.IsPause = false
			l.pauseTime = 0
		}
		return
	}
	l.FindPath()
	next, min := l.p, l.pathMap[l.p]
	for _, i := range l.getOpt().TileAdjacency[l.p] {
		if l.pathMap[i] < min {
			min = l.pathMap[i]
			next = i
		}
	}
	currentP := l.getOpt().IntToTile(l.p)
	nextP := l.getOpt().IntToTile(next)
	var dir string
	if nextP.X < currentP.X {
		dir = "left"
	} else if nextP.X > currentP.X {
		dir = "right"
	} else if nextP.Y < currentP.Y {
		dir = "up"
	} else if nextP.Y > currentP.Y {
		dir = "down"
	} else {
		dir = "center"
	}
	currentRealP := l.getOpt().RealPosition(currentP)
	nextRealP := l.getOpt().RealPosition(nextP)
	speed := l.getOpt().laserSpeed(l.match.Energy)
	delta := speed * dt
	var dir2 string
	dx, dy := l.Pos.X-currentRealP.X, l.Pos.Y-currentRealP.Y
	if math.Abs(dy) > math.Abs(dx) {
		if math.Abs(dy) < delta {
			dir2 = "center"
		} else if dy < 0 {
			dir2 = "up"
		} else {
			dir2 = "down"
		}
	} else {
		if math.Abs(dx) < delta {
			dir2 = "center"
		} else if dx < 0 {
			dir2 = "left"
		} else {
			dir2 = "right"
		}
	}
	var destRealP RP
	if dir == dir2 || dir2 == "center" { // move to next directly
		destRealP = nextRealP
	} else { // move to current first
		destRealP = currentRealP
	}
	dx, dy = destRealP.X-l.Pos.X, destRealP.Y-l.Pos.Y
	pos := l.Pos
	if math.Abs(dx) < math.Abs(dy) {
		if dy > 0 {
			pos.Y = math.Min(destRealP.Y, l.Pos.Y+delta)
		} else {
			pos.Y = math.Max(destRealP.Y, l.Pos.Y-delta)
		}
	} else {
		if dx > 0 {
			pos.X = math.Min(destRealP.X, l.Pos.X+delta)
		} else {
			pos.X = math.Max(destRealP.X, l.Pos.X-delta)
		}
	}
	l.Pos = pos
	newPos, _ := l.getOpt().TilePosition(l.Pos)
	l.p = l.getOpt().TilePosToInt(newPos)
	size := float64(l.getOpt().ArenaCellSize)
	rect := Rect{
		X: pos.X - size/2,
		Y: pos.Y - size/2,
		W: size,
		H: size,
	}
	shouldPause := false
	for _, player := range l.match.Member {
		if player.InvincibleTime > 0 {
			continue
		}
		playerSize := float64(l.getOpt().PlayerSize)
		playerRect := Rect{player.Pos.X - playerSize/2, player.Pos.Y - playerSize/2, playerSize, playerSize}
		if l.getOpt().Collide(&rect, &playerRect) {
			shouldPause = true
			player.InvincibleTime = l.getOpt().PlayerInvincibleTime
			player.HitCount += 1
			var punish int
			if l.match.Mode == "g" {
				punish = int(float64(l.match.Gold) * l.getOpt().Mode1TouchPunish)
			} else {
				punish = l.getOpt().Mode2TouchPunish
			}
			l.match.Gold = MaxInt(l.match.Gold-punish, 0)
			player.Gold -= punish
			player.LostGold += punish
		}
	}
	if shouldPause {
		l.Pause(l.getOpt().LaserPauseTime)
	}
}

func (l *Laser) FindPath() {
	p, _ := l.getOpt().TilePosition(l.player.Pos)
	dest := l.getOpt().TilePosToInt(p)
	if l.dest == dest {
		return
	}
	l.dest = dest
	for i := 0; i < l.getOpt().ArenaWidth*l.getOpt().ArenaHeight; i++ {
		l.pathMap[i] = 10000
	}
	var fill func(x int, v int)
	fill = func(x int, v int) {
		l.pathMap[x] = v
		for _, i := range l.getOpt().TileAdjacency[x] {
			if l.pathMap[i] > v+1 {
				fill(i, v+1)
			}
		}
	}
	fill(l.dest, 0)
}

func (l *Laser) getOpt() *MatchOptions {
	return l.match.opt
}
