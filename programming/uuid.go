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
