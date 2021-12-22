package programming

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-lib/programming"
)

// postUuidOutput is the output of the "POST /programming/uuid" action
type postUuidOutput struct {
	UUID string `json:"uuid"`
}

// SetRouterGroup defines all the routes for the programming functions
func SetRouterGroup(base *gin.RouterGroup) *gin.RouterGroup {
	programmingGroup := base.Group("/programming")
	{
		programmingGroup.POST("/uuid", postUuid())
		// Add here more functions in the programming category
	}

	return programmingGroup
}

// postUuid handles the uuid request.
//
// Reads the "no-hyphens" parameter from the query string to support
// UUIDs without hyphens.
//
// It returns 200 on success.
func postUuid() gin.HandlerFunc {
	return func(c *gin.Context) {
		noHyphensParamValue := c.Query("no-hyphens")
		withoutHyphens := noHyphensParamValue == "true"

		uuid := programming.NewUuid(withoutHyphens)
		output := postUuidOutput{UUID: uuid}

		c.JSON(http.StatusOK, output)
	}
}
