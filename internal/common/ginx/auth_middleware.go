package ginx

import (
	"net/http"
	"strings"

	"godir/internal/common/exterr"
	"godir/internal/common/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, Fail(exterr.Newf(-1, "未提供认证token")))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, Fail(exterr.Newf(-1, "token格式错误")))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, Fail(exterr.Newf(10000001, "无效的token")))
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userId", claims.UserID)
		c.Set("userName", claims.Username)

		c.Set("userInfo", *claims)

		c.Next()
	}
}
