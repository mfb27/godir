package handler

import (
	"godir/internal/common/ginx"
	"godir/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	// g := r.Group("/user")
	// g.POST("/create", ginx.WrapHandlerObj((*user.User).Create))

	// 需要认证的用户路由组
	protected := r.Group("/user")
	protected.Use(ginx.AuthMiddleware())
	{
		protected.GET("/profile", ginx.WrapHandlerObj((*user.User).Profile))
		protected.PUT("/profile", ginx.WrapHandlerObj((*user.User).UpdateProfile))
		protected.POST("/avatar", ginx.WrapHandlerObj((*user.User).UploadAvatar))
	}
}
