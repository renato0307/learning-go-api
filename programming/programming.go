package programming

import (
	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-lib/programming"
	"github.com/rs/zerolog/log"
)

// SetRouterGroup defines all the routes for the programming functions
func SetRouterGroup(p programming.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	log.Debug().Msg("setting router group for: programming")

	programmingGroup := base.Group("/programming")
	{
		programmingGroup.POST("/uuid", postUuid(p))
		programmingGroup.POST("/jwt", postJwtDebugger(p))
		// Add here more functions in the programming category
	}

	return programmingGroup
}
