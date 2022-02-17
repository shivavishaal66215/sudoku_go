package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main(){
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/test",func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, "Hello there!")	
	})
	r.POST("/gen-sudoku", HandleGenSudoku)
	r.POST("/register", HandleRegister)

	r.Run("localhost:8000")
}