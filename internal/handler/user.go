package handler

import (
	"godir/internal/handler/auth"
	"godir/internal/handler/base"
	"godir/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	g := r.Group("/user")
	g.POST("/create", base.WrapHandlerObj((*user.User).Create))

	// 需要认证的用户路由组
	protected := r.Group("/user")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/profile", base.WrapHandlerObj((*user.User).Profile))
		protected.PUT("/profile", base.WrapHandlerObj((*user.User).UpdateProfile))
		protected.POST("/avatar", base.WrapHandlerObj((*user.User).UploadAvatar))
	}
}
