package finance

import (
	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-lib/finance"
)

// SetRouterGroup defines all the routes for the finance functions
func SetRouterGroup(f finance.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	financeGroup := base.Group("/finance")
	{
		financeGroup.GET("/currconv", getCurrConv(f))
		// Add here more functions in the finance category
	}

	return financeGroup
}
