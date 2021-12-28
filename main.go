package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-api/finance" // new
	"github.com/renato0307/learning-go-api/programming"
	financelib "github.com/renato0307/learning-go-lib/finance" // new
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

	useDefaultUrl := ""                                        // new
	apiKey := getRequiredEnv("CURRCONV_API_KEY")               // new
	f := financelib.NewFinanceFunctions(useDefaultUrl, apiKey) // new
	finance.SetRouterGroup(&f, base)                           // new

	r.Run()
}

func getRequiredEnv(key string) string { // new
	value, exists := os.LookupEnv(key)

	if !exists {
		panic(fmt.Sprintf("error: %s environment variable was not defined", key))
	}

	return value
}
