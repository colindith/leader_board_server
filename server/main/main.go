package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/colindith/leader_board_server/server/handler"
)

func main() {
	r := gin.Default()
	r.GET("/api/ping", handler.Ping)

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Print("[ERROR] server ends with error: " + err.Error())
	}
}