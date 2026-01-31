package ai

import (
	"fmt"

	"godir/internal/common/ginx"
	"godir/internal/model"
	"godir/internal/types"

	"github.com/gin-gonic/gin"
)

type Ai struct {
	ginx.BaseHandler
}

func (h *Ai) New() ginx.BaseHandlerInterface {
	return new(Ai)
}

func (h *Ai) List(c *gin.Context, req *types.AiAppListReq) (*types.AiAppListResp, error) {
	var apps []model.GodirAiApp
	if err := h.DB.Order("created_at DESC").Find(&apps).Error; err != nil {
		return nil, fmt.Errorf("查询AI应用列表失败: %w", err)
	}

	appList := make([]types.AiAppInfo, 0, len(apps))
	for _, app := range apps {
		appList = append(appList, types.AiAppInfo{
			ID:    app.ID,
			Name:  app.Name,
			AppID: app.AppID,
			Desc:  app.Desc,
			Icon:  app.Icon,
			Cover: app.Cover,
		})
	}

	return &types.AiAppListResp{
		Apps: appList,
	}, nil
}

