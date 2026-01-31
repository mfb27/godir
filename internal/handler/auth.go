package handler

import (
	"godir/internal/common/ginx"
	"godir/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRouter(r *gin.Engine) {
	// 不需要认证的路由组
	public := r.Group("/auth")
	{
		public.POST("/register", ginx.WrapHandlerObj((*auth.Auth).Register))
		public.POST("/login", ginx.WrapHandlerObj((*auth.Auth).Login))
	}

	// 需要认证的路由组
	protected := r.Group("/auth")
	protected.Use(ginx.AuthMiddleware())
	{
		protected.POST("/logout", ginx.WrapHandlerObj((*auth.Auth).Logout))
	}
}
