package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-api/internal/apierror"
	"github.com/rs/zerolog/log"
)

func Authorizer() gin.HandlerFunc {
	return func(c *gin.Context) {

		// extracts the required scope from the URL
		p1 := c.Request.URL.Path               // p1 = "/v1/programming/uuid"
		p2 := strings.ReplaceAll(p1, "/", "-") // p2 = "-v1-programming-uuid"
		p3 := strings.Split(p2, "-")           // p3 = ["", "v1", "programming", "uuid"]
		p4 := p3[2:]                           // p4 = ["programming", "uuid"]
		scope := strings.Join(p4, "-")         // scope = "programming-uuid"
		log.Debug().Msgf("scope for url is %s", scope)

		// gets the scopes added by the Authenticator middleware
		clientScopes := c.GetString(ScopeKey)
		clientScopesList := strings.Split(clientScopes, " ")
		log.Debug().Msgf("client scope list is %s", clientScopes)

		// tries to find a matching scope
		// scopes have the following format
		// https://learninggolang.com/programming-jwtdebugger
		found := false
		for _, clientScope := range clientScopesList {
			found = strings.HasSuffix(clientScope, scope)
			if found {
				break
			}
		}

		// returns forbidden (HTTP status 403) if no valid scope is found
		if !found {
			log.Debug().Msg("no scope found for current route")
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				apierror.New("Forbidden"))
			return
		}

		log.Debug().Msg("valid client scope found for current route")
	}
}
