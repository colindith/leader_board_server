package service

import (
	"log"
	"time"

	"github.com/colindith/leader_board_server/server/remote"
)

func UpdateScore(clientID int64, score int32) (errorCode int32) {
	scoreObj := &remote.Score{
		ClientID : clientID,
		Score    : score,
		Timestamp: time.Now().Unix(),
	}
	code := remote.GetScoreDB().CreateScore(scoreObj)
	if code != remote.DB_SUCCESS {
		log.Printf("[WARN] update score DB failed, error code: %v", code)
	}
	return code
}
