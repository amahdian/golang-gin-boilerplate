package middleware

import (
	"net/http"

	"github.com/amahdian/golang-gin-boilerplate/domain/contracts/resp"
	"github.com/amahdian/golang-gin-boilerplate/svc/auth"

	"github.com/gin-gonic/gin"
)

func VerifyAuth(authenticator auth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		ctx, err := authenticator.Verify(r)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.NewErrorResponse(err))
			return
		}
		c.Request = r.WithContext(ctx)
		c.Next()
	}
}
