package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tipbk/blog-backend/config"
	"github.com/tipbk/blog-backend/service"
	"github.com/tipbk/blog-backend/util"
)

type AuthMiddleware interface {
	AuthAccessTokenMiddleware(c *gin.Context)
}

type authMiddleware struct {
	envConfig   *config.EnvConfig
	userService service.UserService
}

func NewAuthMiddleware(envConfig *config.EnvConfig, userService service.UserService) AuthMiddleware {
	return &authMiddleware{
		envConfig:   envConfig,
		userService: userService,
	}
}

func (m *authMiddleware) AuthAccessTokenMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	tokenArr := strings.Split(authHeader, "Bearer ")
	// invalid token
	if len(tokenArr) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, util.GenerateFailedResponse("token is invalid"))
		return
	}
	jwt, err := util.ValidateAccessToken(m.envConfig.AccessTokenSecret, tokenArr[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, util.GenerateFailedResponse("token is invalid"))
		return
	}

	userID := jwt["userID"].(string)
	user, err := m.userService.FindUserWithUserID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, util.GenerateFailedResponse(err.Error()))
		return
	}

	// set user
	c.Set("user", user)

	c.Next()
}
