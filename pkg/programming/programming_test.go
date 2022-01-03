package programming

import (
	"github.com/gin-gonic/gin"
	programminglib "github.com/renato0307/learning-go-lib/programming"
)

func setupGin(mockInterface *programminglib.MockInterface) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(mockInterface, v1)

	return r
}
