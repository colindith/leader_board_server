package remote

import (
	"gorm.io/gorm"
)

type Score struct {
	gorm.Model
	ClientID  int64 `json:"client_id"`
	Score     int32 `json:"score"`
	Timestamp int64 `json:"timestamp"`
}

type ScoreDB interface {
	GetTop10() (int32, map[int]interface{})
	CreateScore(score *Score) int32
}

type ScoreDBImp struct {}

var scoreDBObj ScoreDB = &ScoreDBImp{}

func GetScoreDB() ScoreDB {
	return scoreDBObj
}

func SetScoreDB(s ScoreDB) {
	scoreDBObj = s
}

func (s *ScoreDBImp) GetTop10() (int32, map[int]interface{}) {
	results := map[int]interface{}{}
	GetDB().
		Limit(10).
		Distinct("client_id").
		Order("score , timestamp desc").
		Find(&results)
	return 0, results
}

func (s *ScoreDBImp) CreateScore(score *Score) int32 {
	result := GetDB().Create(&score)
	if result.Error != nil {
		return 1
	}
	return 0
}