package remote

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ScoreLeaderBoardCacheKey     = "score_leader_board"
	ScoreLeaderBoardCacheTimeout = 600 // seconds
)

type Score struct {
	gorm.Model
	ClientID  string  `json:"client_id" gorm:"type:varchar(32);unique_index:idx_client_id_group"`
	Score     float64 `json:"score" gorm:"type:decimal(10,2);index:idx_score_group,priority:2"`
	Group     int64   `json:"group" gorm:"unique_index:idx_client_id_group;index:idx_score_group,priority:1"`
	Timestamp int64   `json:"timestamp"`
}

type ScoreDB interface {
	GetTop10Score(group int64) (int32, []*Score)
	UpdateScore(score *Score) int32
}

type ScoreDBImp struct {}

var scoreDBObj ScoreDB = &ScoreDBImp{}

func GetScoreDB() ScoreDB {
	return scoreDBObj
}

func SetScoreDB(s ScoreDB) {
	scoreDBObj = s
}

// GetTop10Score query the 10 highest scores from the current group
func (s *ScoreDBImp) GetTop10Score(group int64) (int32, []*Score) {
	scores := make([]*Score, 0, 10)
	GetDB().
		Limit(10).
		Where(&Score{Group: group}).
		Order("score desc").
		Find(&scores)
	if len(scores) == 0 {
		return DB_EMPTY_RESULT_ERROR, nil
	}
	return CODE_SUCCESS, scores
}

func (s *ScoreDBImp) UpdateScore(score *Score) int32 {
	tx := GetDB().Begin()
	oldScore := Score{}
	tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&Score{ClientID: score.ClientID, Group: score.Group}).
		First(&oldScore)

	if oldScore.ID == 0 {
		// client_id not exist in DB, create new record for it
		tx.Create(&score)
		if err := tx.Commit().Error; err != nil {
			log.Printf("[ERROR] create new score error: %v", err)
			return DB_ERROR
		}
		return CODE_SUCCESS
	}
	if score.Score > oldScore.Score {
		// update the record
		tx.
			Model(&Score{}).
			Where(&Score{ClientID: score.ClientID, Group: score.Group}).
			Update("score", score.Score)
		if err := tx.Commit().Error; err != nil {
			log.Printf("[ERROR] update score error: %v, clientID: %v", err, score.ClientID)
			return DB_ERROR
		}
		return CODE_SUCCESS
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("[ERROR] release lock error: %v", err)
		return DB_ERROR
	}
	return DB_NOT_UPDATED
}