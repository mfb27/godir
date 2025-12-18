package ginx

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(log *zap.SugaredLogger) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源，生产环境建议指定具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: false, // 当AllowOrigins为"*"时，AllowCredentials必须为false
		MaxAge:           12 * time.Hour,
	}))

	r.Use(RequestIDLoggerMiddleware(log))
	r.Use(gin.Recovery())

	return r
}

// RequestIDLoggerMiddleware 生成/透传请求ID，创建带 request_id 的子日志并注入上下文
func RequestIDLoggerMiddleware(base *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = generateRequestID()
		}
		c.Writer.Header().Set("X-Request-ID", rid)

		child := base.With(zap.String("request_id", rid))

		// 注入到 request context 和 gin.Keys，便于下游获取
		// ctx := logger.WithContext(c.Request.Context(), child)
		// c.Request = c.Request.WithContext(ctx)
		c.Set("logger", child)

		c.Next()
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
