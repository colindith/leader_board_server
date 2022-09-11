package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/colindith/leader_board_server/server/remote"
	"github.com/colindith/leader_board_server/server/service"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func UpdateScore(c *gin.Context) {
	clientID := c.GetHeader("ClientId")

	json := service.UpdateScoreReq{}
	err := c.BindJSON(&json)
	if err != nil {
		log.Printf("[ERROR] error request payload %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "not ok",
			"message": "error request payload",
		})
		return
	}
	log.Print("[DEBUG] update_score: ", &json)
	code := service.UpdateScore(clientID, json.Score)
	if code != remote.DB_SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status": "not ok",
			"message": fmt.Sprintf("update database failed: %v", code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func GetTop10Score(c *gin.Context) {
	resp := &service.GetTop10ScoreResp{}
	code := service.GetTop10Score(resp)
	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "not ok",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"topPlayers": resp,
	})
}
