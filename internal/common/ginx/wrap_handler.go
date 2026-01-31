package ginx

import (
	"context"
	"godir/internal/common/exterr"
	"godir/internal/common/logger"
	"godir/internal/common/svc"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseHandlerInterface interface {
	New() BaseHandlerInterface
	SetLogger(l *zap.SugaredLogger)
	SetCtx(ctx context.Context)
	SetRequestID(id string)
	SetDB(db *gorm.DB)
	SetSvc(svc *svc.ServiceContext)
}

// BaseHandler 提供通用的上下文、配置与依赖注入能力
type BaseHandler struct {
	Log       *zap.SugaredLogger
	Ctx       context.Context
	RequestID string
	DB        *gorm.DB
	Svc       *svc.ServiceContext
}

func (b *BaseHandler) SetLogger(l *zap.SugaredLogger) { b.Log = l }
func (b *BaseHandler) SetCtx(ctx context.Context)     { b.Ctx = ctx }
func (b *BaseHandler) SetRequestID(id string)         { b.RequestID = id }
func (b *BaseHandler) SetDB(db *gorm.DB)              { b.DB = db }
func (b *BaseHandler) SetSvc(svc *svc.ServiceContext) { b.Svc = svc }

func WrapHandlerObj[T BaseHandlerInterface, X any, Y any](method func(T, *gin.Context, *X) (*Y, error)) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req X
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, Fail(exterr.Newf(-1, "参数错误: %v", err)))
			c.Abort()
			return
		}

		var t T
		obj := t.New().(T)
		obj.SetLogger(logger.FromContext(c))
		obj.SetCtx(c.Request.Context())
		obj.SetRequestID(c.GetHeader("X-Request-ID"))
		obj.SetDB(svc.DB())
		obj.SetSvc(svc.Get())

		resp, err := method(obj, c, &req)
		if err != nil {
			c.JSON(http.StatusOK, Fail(err))
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, Success(resp))
	}
}
