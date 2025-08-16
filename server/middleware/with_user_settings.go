package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/amahdian/golang-gin-boilerplate/server/utils"
	"github.com/amahdian/golang-gin-boilerplate/storage"
	"github.com/amahdian/golang-gin-boilerplate/svc/auth"
	"github.com/gin-gonic/gin"
)

func WithUserSettings(stg storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userInfo := auth.UserInfoFromCtx(ctx)
		usStg := stg.User(ctx)
		userSettings, err := usStg.FindByEmail(userInfo.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("could not fetch user settings: %v", err),
			})
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), utils.UserSettingsContextKey, userSettings))
		c.Next()
	}
}
