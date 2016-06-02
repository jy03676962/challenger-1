package core

type LaserInterface interface {
	Pause(t float64)
	IsFollow(cid string) bool
	Tick(dt float64)
}
