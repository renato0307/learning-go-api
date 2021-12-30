package programming

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-lib/programming"
)

// postJwtDebuggerOutput is the output of the "POST /programming/jwt" action
type postJwtDebuggerOutput struct {
	Header  string `json:"header"`
	Payload string `json:"payload"`
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
