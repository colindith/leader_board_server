package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/colindith/leader_board_server/server/handler"
	"github.com/colindith/leader_board_server/server/remote"
)

func main() {
	r := gin.Default()
	err := remote.InitDB()
	if err != nil {
		return
	}
	err = remote.GetDB().AutoMigrate(&remote.Score{})
	if err != nil {
		log.Printf("[ERROR] gorm migrate error: %v", err)
		return
	}

	r.GET("/api/ping", handler.Ping)
	r.POST("/api/score", handler.UpdateScore)
	r.GET("/api/top10_score", handler.GetTop10Score)

	err = r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Print("[ERROR] server ends with error: " + err.Error())
	}
}