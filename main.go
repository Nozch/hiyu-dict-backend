package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
type Metaphor struct {
	ID		int `json:"id"`
	Adjective string `json:"adjective"`
	Metaphor 	string `json:"metaphor"`
}
func main() {
	router := gin.Default()

	

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	metaphors := []Metaphor{
		{ID: 1, Adjective: "熱い", Metaphor: "中華あん"},
		{ID: 2, Adjective: "てらてら", Metaphor: "油まみれ"},
		{ID: 3, Adjective: "芳しい", Metaphor: "チーズ"},
	}

	router.GET("/metaphors", func(c *gin.Context) {
		c.JSON(http.StatusOK, metaphors)
	})


	router.Run(":3000")
}