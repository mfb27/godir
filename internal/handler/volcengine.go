package handler

import (
	"godir/internal/common/ginx"
	"godir/internal/handler/volcengine"

	"github.com/gin-gonic/gin"
)

func RegisterVolcEngineRouter(r *gin.Engine) {
	protected := r.Group("/volcengine")
	protected.Use(ginx.AuthMiddleware())
	{
		protected.POST("/knowledge-base/create", ginx.WrapHandlerObj((*volcengine.VolcEngine).CreateKnowledgeBase))
		protected.GET("/knowledge-base/list", ginx.WrapHandlerObj((*volcengine.VolcEngine).ListKnowledgeBase))
		protected.POST("/knowledge-base/delete", ginx.WrapHandlerObj((*volcengine.VolcEngine).DeleteKnowledgeBase))
		protected.POST("/document/upload", ginx.WrapHandlerObj((*volcengine.VolcEngine).UploadDocument))
		protected.GET("/document/list", ginx.WrapHandlerObj((*volcengine.VolcEngine).ListDocument))
		protected.POST("/document/delete", ginx.WrapHandlerObj((*volcengine.VolcEngine).DeleteDocument))
		protected.POST("/chat", ginx.WrapHandlerObj((*volcengine.VolcEngine).Chat))
		protected.POST("/search", ginx.WrapHandlerObj((*volcengine.VolcEngine).Search))
	}
}
