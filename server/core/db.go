package core

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type PlayerData struct {
	gorm.Model
	Name      string  `json:"name"`
	Gold      int     `json:"gold"`
	LostGold  int     `json:"lostGold"`
	Energy    float64 `json:"energy"`
	Combo     int     `json:"combo"`
	Grade     string  `json:"grade"`
	Level     int     `json:"level"`
	LevelData string  `json:"levelData"`
	HitCount  int     `json:"hitCount"`
}

type MatchData struct {
	gorm.Model
	Mode         string       `json:"mode"`
	Elasped      float64      `json:"elasped"`
	Gold         int          `json:"gold"`
	Member       []PlayerData `json:"member"`
	RampageCount int          `json:"rampageCount"`
}

type DB struct {
	conn *gorm.DB
}

func NewDb() *DB {
	return &DB{}
}

func (db *DB) connect(path string) error {
	conn, err := gorm.Open("sqlite3", path)
	conn.AutoMigrate(&MatchData{}, &PlayerData{})
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *DB) saveMatch(d *MatchData) uint {
	db.conn.Create(d)
	return d.ID
}

func (db *DB) updateMatchData(m *MatchData) {
	db.conn.Save(&m)
}

func (db *DB) getLatestMatch() *MatchData {
	var m MatchData
	db.conn.Last(&m)
	return &m
}
