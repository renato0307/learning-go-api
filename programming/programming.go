package programming

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-lib/programming"
)

// postUuidOutput is the output of the "POST /programming/uuid" action
type postUuidOutput struct {
	UUID string `json:"uuid"`
}

// postJwtDebuggerOutput is the output of the "POST /programming/jwt" action
type postJwtDebuggerOutput struct {
	Header  string `json:"header"`
	Payload string `json:"payload"`
}

// SetRouterGroup defines all the routes for the programming functions
func SetRouterGroup(p programming.Interface, base *gin.RouterGroup) *gin.RouterGroup {
	programmingGroup := base.Group("/programming")
	{
		programmingGroup.POST("/uuid", postUuid(p))
		programmingGroup.POST("/jwt", postJwtDebugger(p))
		// Add here more functions in the programming category
	}

	return programmingGroup
}

// postUuid handles the uuid request.
//
// Reads the "no-hyphens" parameter from the query string to support
// UUIDs without hyphens.
//
// It returns HTTP 200 on success.
func postUuid(p programming.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		noHyphensParamValue := c.Query("no-hyphens")
		withoutHyphens := noHyphensParamValue == "true"

		uuid := p.NewUuid(withoutHyphens)
		output := postUuidOutput{UUID: uuid}

		c.JSON(http.StatusOK, output)
	}
}

// postJwtDebugger handles the JWT debug request.
//
// It returns HTTP 200 on success.
// Returns HTTP 400 if the token is not valid.
func postJwtDebugger(p programming.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, "error reading body")
			return
		}

		tokenString := string(tokenBytes)
		header, payload, err := p.DebugJWT(tokenString)
		if err != nil {
			message := fmt.Sprintf("invalid token: %s", err.Error())
			c.JSON(http.StatusBadRequest, message)
			return
		}

		output := postJwtDebuggerOutput{
			Header:  header,
			Payload: payload,
		}
		c.JSON(http.StatusOK, output)
	}
}
