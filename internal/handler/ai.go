package handler

import (
	"godir/internal/common/ginx"
	"godir/internal/handler/ai"

	"github.com/gin-gonic/gin"
)

func RegisterAiRouter(r *gin.Engine) {
	public := r.Group("/ai")
	{
		public.GET("/apps", ginx.WrapHandlerObj((*ai.Ai).List))
	}
}

