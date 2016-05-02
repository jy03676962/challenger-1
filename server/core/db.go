package core

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type PlayerData struct {
	gorm.Model
	Name     string
	Gold     int
	LostGold int
	Energy   float64
	Combo    int
	Grade    string
	Level    int
}

type MatchData struct {
	gorm.Model
	Mode    string
	Elasped float64
	Gold    int
	Member  []PlayerData
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
