package handler

import (
	"godir/internal/handler/auth"

	"godir/internal/handler/base"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRouter(r *gin.Engine) {
	// 不需要认证的路由组
	public := r.Group("/auth")
	{
		public.POST("/register", base.WrapHandlerObj((*auth.Auth).Register))
		public.POST("/login", base.WrapHandlerObj((*auth.Auth).Login))
	}

	// 需要认证的路由组
	protected := r.Group("/auth")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/logout", base.WrapHandlerObj((*auth.Auth).Logout))
	}
}
