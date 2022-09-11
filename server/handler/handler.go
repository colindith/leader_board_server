package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/colindith/leader_board_server/server/service"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func UpdateScore(c *gin.Context) {
	json := service.UpdateScoreReq{}
	err := c.BindJSON(&json)
	if err != nil {
		log.Printf("[ERROR] error request payload %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error request payload",
			"success": "false",
		})
		return
	}
	log.Print("[DEBUG] update_score: ", &json)
	code := service.UpdateScore(json.ClientID, json.Score)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"success": "true",
	})
}

func GetTop10Score(c *gin.Context) {
	resp := &service.GetTop10ScoreResp{}
	code := service.GetTop10Score(resp)
	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"success": "false",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"success": "true",
		"score_list": resp,
	})
}
