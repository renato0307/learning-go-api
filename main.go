package main

import (
	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-api/programming"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, welcome to the learning-go-api",
		})
	})

	base := r.Group("/v1")
	programming.SetRouterGroup(base)

	r.Run()
}
