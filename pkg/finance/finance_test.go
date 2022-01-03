package finance

import (
	"github.com/gin-gonic/gin"
	financelib "github.com/renato0307/learning-go-lib/finance"
)

func setupGin(mockInterface *financelib.MockInterface) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(mockInterface, v1)

	return r
}
