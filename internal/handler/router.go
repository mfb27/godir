package handler

import (
	"godir/internal/common/redis"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	RegisterAuthRouter(r)
	RegisterUserRouter(r)
	RegisterMaterialRouter(r)
	RegisterVolcEngineRouter(r)

	redis.StartThumbnailWorker()
}
