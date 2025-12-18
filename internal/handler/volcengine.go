package handler

import (
	"godir/internal/handler/auth"
	"godir/internal/handler/base"
	"godir/internal/handler/volcengine"

	"github.com/gin-gonic/gin"
)

func RegisterVolcEngineRouter(r *gin.Engine) {
	protected := r.Group("/volcengine")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/knowledge-base/create", base.WrapHandlerObj((*volcengine.VolcEngine).CreateKnowledgeBase))
		protected.GET("/knowledge-base/list", base.WrapHandlerObj((*volcengine.VolcEngine).ListKnowledgeBase))
		protected.POST("/knowledge-base/delete", base.WrapHandlerObj((*volcengine.VolcEngine).DeleteKnowledgeBase))
		protected.POST("/document/upload", base.WrapHandlerObj((*volcengine.VolcEngine).UploadDocument))
		protected.GET("/document/list", base.WrapHandlerObj((*volcengine.VolcEngine).ListDocument))
		protected.POST("/document/delete", base.WrapHandlerObj((*volcengine.VolcEngine).DeleteDocument))
		protected.POST("/chat", base.WrapHandlerObj((*volcengine.VolcEngine).Chat))
		protected.POST("/search", base.WrapHandlerObj((*volcengine.VolcEngine).Search))
	}
}
