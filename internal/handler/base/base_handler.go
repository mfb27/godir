package base

import (
	"context"

	"godir/internal/common/svc"

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
