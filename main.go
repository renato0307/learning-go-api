package main

import (
	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-api/programming"
	programminglib "github.com/renato0307/learning-go-lib/programming"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, welcome to the learning-go-api",
		})
	})

	base := r.Group("/v1")

	p := programminglib.ProgrammingFunctions{}
	programming.SetRouterGroup(&p, base)

	r.Run()
}
