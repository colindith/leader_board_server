package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func UpdateScore(c *gin.Context) {
	json := Score{}
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
	c.JSON(http.StatusOK, gin.H{
		"message": &json,
		"success": "true",
	})
}

func GetScore(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type Score struct {
	ClientID int `json:"client_id"`
	Score    int `json:"score"`
}