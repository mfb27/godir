package handler

import (
	"godir/internal/handler/auth"
	"godir/internal/handler/base"
	"godir/internal/handler/material"

	"github.com/gin-gonic/gin"
)

func RegisterMaterialRouter(r *gin.Engine) {
	// 需要认证的路由组
	protected := r.Group("/material")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/upload-token", base.WrapHandlerObj((*material.Material).GetUploadToken))
		protected.POST("/save", base.WrapHandlerObj((*material.Material).Save))
		protected.GET("/list", base.WrapHandlerObj((*material.Material).List))
		protected.GET("/search", base.WrapHandlerObj((*material.Material).Search))
		protected.POST("/delete", base.WrapHandlerObj((*material.Material).BatchDelete))
		protected.POST("/update-name", base.WrapHandlerObj((*material.Material).UpdateMaterialName))
		protected.POST("/publish", base.WrapHandlerObj((*material.Material).Publish))
		protected.POST("/published/like", base.WrapHandlerObj((*material.Material).LikePublish))
		protected.POST("/published/unlike", base.WrapHandlerObj((*material.Material).UnlikePublish))
	}

	// 公开的路由组（无需认证）
	public := r.Group("/public")
	{
		public.GET("/published", base.WrapHandlerObj((*material.Material).ListPublished))
	}
}
