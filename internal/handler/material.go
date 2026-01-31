package handler

import (
	"godir/internal/common/ginx"
	"godir/internal/handler/material"

	"github.com/gin-gonic/gin"
)

func RegisterMaterialRouter(r *gin.Engine) {
	// 需要认证的路由组
	protected := r.Group("/material")
	protected.Use(ginx.AuthMiddleware())
	{
		protected.POST("/upload-token", ginx.WrapHandlerObj((*material.Material).GetUploadToken))
		protected.POST("/save", ginx.WrapHandlerObj((*material.Material).Save))
		protected.GET("/list", ginx.WrapHandlerObj((*material.Material).List))
		protected.GET("/search", ginx.WrapHandlerObj((*material.Material).Search))
		protected.POST("/delete", ginx.WrapHandlerObj((*material.Material).BatchDelete))
		protected.POST("/update-name", ginx.WrapHandlerObj((*material.Material).UpdateMaterialName))
		protected.POST("/publish", ginx.WrapHandlerObj((*material.Material).Publish))
		protected.POST("/published/like", ginx.WrapHandlerObj((*material.Material).LikePublish))
		protected.POST("/published/unlike", ginx.WrapHandlerObj((*material.Material).UnlikePublish))
	}

	// 公开的路由组（无需认证）
	public := r.Group("/public")
	{
		public.GET("/published", ginx.WrapHandlerObj((*material.Material).ListPublished))
	}
}
