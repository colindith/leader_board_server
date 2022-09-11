package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/colindith/leader_board_server/server/remote"
)

func UpdateScore(clientID string, score float64) (code int32) {
	ts := time.Now().Unix()
	scoreObj := &remote.Score{
		ClientID : clientID,
		Score    : score,
		Timestamp: ts,
		Group    : getGroupFromTimestamp(ts),
	}
	code = remote.GetScoreDB().UpdateScore(scoreObj)
	if code != remote.CODE_SUCCESS && code != remote.DB_NOT_UPDATED {
		log.Printf("[WARN] update score DB failed, error code: %v", code)
		return
	}
	if code == remote.DB_NOT_UPDATED {
		return remote.CODE_SUCCESS
	}

	// delete the cache if the ClientID is in the current leaderboard
	top10ScoreResp, _ := getTop10ScoreCache()
	if top10ScoreResp == nil {
		// not cached
		return
	}
	if len(top10ScoreResp.ClientIDScoreList) < 10 {
		remote.GetRedis().Delete(remote.ScoreLeaderBoardCacheKey)
		return
	}
	for _, clientIDScore := range top10ScoreResp.ClientIDScoreList {
		if clientIDScore.ClientID == clientID {
			// delete cache
			remote.GetRedis().Delete(remote.ScoreLeaderBoardCacheKey)
			break
		}
	}
	return
}

type UpdateScoreReq struct {
	Score    float64 `json:"score"`
}

func GetTop10Score(resp *GetTop10ScoreResp) (errCode int32) {
	cachedResp, _ := getTop10ScoreCache()
	if cachedResp != nil && isCurGroup(cachedResp.Group) {
		log.Print("[DEBUG] fetch leaderboard from cache")
		resp.Group = cachedResp.Group
		resp.ClientIDScoreList = cachedResp.ClientIDScoreList
		return remote.CODE_SUCCESS
	}

	curGroup := getGroupFromTimestamp(time.Now().Unix())

	code, rawScores := remote.GetScoreDB().GetTop10Score(curGroup)
	log.Print("[DEBUG] top 10 scores code, result: ", code, rawScores)
	if code != remote.CODE_SUCCESS && code != remote.DB_EMPTY_RESULT_ERROR {
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
	resp.Group = curGroup
	// update the cache
	if err := updateTop10ScoreCache(resp); err != nil {
		log.Printf("[ERROR] update leaderboard cache failed: %v", err)
	}

	return remote.CODE_SUCCESS
}

type GetTop10ScoreResp struct {
	ClientIDScoreList []*ClientIDScore `json:"client_id_score_list"`
	Group              int64           `json:"group"`
}

type ClientIDScore struct {
	ClientID  string  `json:"clientId"`
	Score     float64 `json:"score"`
}

func getGroupFromTimestamp(ts int64) int64 {
	// each group ends in 10 mins, which is equal to 600 secs
	return ts / 600
}

func isCurGroup(group int64) bool {
	curGroup := getGroupFromTimestamp(time.Now().Unix())
	return curGroup == group
}

func getTop10ScoreCache() (res *GetTop10ScoreResp, err error) {
	err, rawData := remote.GetRedis().Get(remote.ScoreLeaderBoardCacheKey)
	if err != nil {
		log.Printf("[WARN] get top10ScoreResp cache error: %v", err)
		return nil, err
	}
	err = json.Unmarshal(rawData, &res)
	if err != nil {
		log.Printf("[ERROR] unmarshal top10ScoreResp error: %v", err)
		return nil, err
	}
	return res, nil
}

func updateTop10ScoreCache(top10ScoreResp *GetTop10ScoreResp) (err error) {

	byteArr, err := json.Marshal(top10ScoreResp)
	if err != nil {
		log.Printf("[ERROR] marshal top10ScoreResp error: %v", err)
		return err
	}

	err = remote.GetRedis().Set(remote.ScoreLeaderBoardCacheKey, byteArr, remote.ScoreLeaderBoardCacheTimeout)
	if err != nil {
		log.Printf("[ERROR] set clientIDScoreList cache error: %v", err)
		return err
	}
	return nil
}