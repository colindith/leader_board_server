package service

import (
	"log"
	"time"

	"github.com/colindith/leader_board_server/server/remote"
)

func UpdateScore(clientID string, score float64) (errorCode int32) {
	ts := time.Now().Unix()
	scoreObj := &remote.Score{
		ClientID : clientID,
		Score    : score,
		Timestamp: ts,
		Group    : getGroupFromTimestamp(ts),
	}
	code := remote.GetScoreDB().UpdateScore(scoreObj)
	if code != remote.DB_SUCCESS {
		log.Printf("[WARN] update score DB failed, error code: %v", code)
	}
	return code
}

type UpdateScoreReq struct {
	ClientID string  `json:"client_id"`
	Score    float64 `json:"score"`
}

func GetTop10Score(resp *GetTop10ScoreResp) (errCode int32) {
	code, rawScores := remote.GetScoreDB().GetTop10Score(getGroupFromTimestamp(time.Now().Unix()))
	log.Print("[DEBUG] top 10 scores code, result: ", code, rawScores)
	if code != remote.DB_SUCCESS {
		return code
	}
	clientIDScoreList := make([]*ClientIDScore, 0, len(rawScores))
	for _, rawScore := range rawScores {
		clientIDScoreList = append(clientIDScoreList, &ClientIDScore{
			ClientID:  rawScore.ClientID,
			Score:     rawScore.Score,
		})
	}
	resp.ClientIDScoreList = clientIDScoreList
	return code
}

type GetTop10ScoreResp struct {
	ClientIDScoreList []*ClientIDScore `json:"client_id_score_list"`
}

type ClientIDScore struct {
	ClientID  string  `json:"clientId"`
	Score     float64 `json:"score"`
}

func getGroupFromTimestamp(ts int64) int64 {
	// each group ends in 10 mins, which is equal to 600 secs
	return ts / 600
}