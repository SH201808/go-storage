package token

import (
	response "file-server/models/Response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, response.Success("未携带token"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, response.Success("请求头auth格式有误"))
			c.Abort()
			return
		}

		mc, err := Parse(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, response.Success("token无效"))
			c.Abort()
			return
		}

		c.Set("userId", int(mc.Data.(float64)))
		c.Next()
	}
}
