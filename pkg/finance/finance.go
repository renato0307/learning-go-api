package finance

import (
	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-lib/finance"

	"github.com/rs/zerolog/log"
)

// SetRouterGroup defines all the routes for the finance functions
func SetRouterGroup(f finance.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	log.Debug().Msg("setting router group for: finance")

	financeGroup := base.Group("/finance")
	{
		financeGroup.GET("/currconv", getCurrConv(f))
		// Add here more functions in the finance category
	}

	return financeGroup
}
