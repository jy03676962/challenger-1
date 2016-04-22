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
	Mode    int
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

func (db *DB) Connect() error {
	conn, err := gorm.Open("sqlite3", "./challenger.db")
	conn.AutoMigrate(&MatchData{}, &PlayerData{})
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *DB) saveMatch(m *Match) *MatchData {
	data := MatchData{}
	data.Mode = m.Mode
	data.Gold = m.Gold
	data.Elasped = m.Elasped
	data.Member = make([]PlayerData, 0)
	for _, player := range m.Member {
		playerData := PlayerData{}
		playerData.Gold = player.Gold
		playerData.Energy = player.Energy
		playerData.LostGold = player.LostGold
		playerData.Combo = player.ComboCount
		data.Member = append(data.Member, playerData)
	}
	db.conn.Create(&data)
	return &data
}

func (db *DB) updateMatchData(m *MatchData) {
	db.conn.Save(&m)
}

func (db *DB) getLatestMatch() *MatchData {
	var m MatchData
	db.conn.Last(&m)
	return &m
}
