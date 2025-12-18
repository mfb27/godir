package base

import (
	"godir/internal/common/exterr"
	"godir/internal/common/logger"
	"godir/internal/common/svc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int64  `json:"code"` // 0成功 其他失败
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) Response {
	return Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
}

// Fail 导出fail函数供其他包使用
func Fail(err error) Response {
	return fail(err)
}

func fail(err error) Response {
	return Response{
		Code: exterr.Code(err),
		Msg:  exterr.Msg(err),
	}
}

// wrapHandlerObj 支持先构造对象实例，再调用其方法。
// 典型用法：wrapHandlerObj(user.New, (*user.User).Create)
// type HasLogger interface{ SetLogger(*zap.SugaredLogger) }
// type HasCtx interface{ SetCtx(ctx context.Context) }
// type HasRequestID interface{ SetRequestID(string) }
// type HasDB interface{ SetDB(db *gorm.DB) }
// type HasCache interface{ SetCache(cache any) }
// type HasConfig interface{ SetConfig(cfg *svc.Config) }
// type HasSvc interface{ SetSvc(svc *svc.ServiceContext) }

func WrapHandlerObj[T BaseHandlerInterface, X any, Y any](method func(T, *gin.Context, *X) (*Y, error)) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req X
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, fail(exterr.Newf(-1, "参数错误: %v", err)))
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
			c.JSON(http.StatusOK, fail(err))
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, Success(resp))
	}
}
