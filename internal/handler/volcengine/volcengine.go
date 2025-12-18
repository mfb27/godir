package volcengine

import (
	"godir/internal/common/exterr"
	"godir/internal/common/svc"
	"godir/internal/common/volcengine"
	"godir/internal/handler/base"
	"godir/internal/types"

	"github.com/gin-gonic/gin"
)

type VolcEngine struct {
	base.BaseHandler
	client *volcengine.Client
}

func (h *VolcEngine) New() base.BaseHandlerInterface {
	cfg := svc.Cfg()
	client := volcengine.NewClient(
		cfg.VolcEngine.AccessKeyID,
		cfg.VolcEngine.SecretAccessKey,
		cfg.VolcEngine.Region,
		cfg.VolcEngine.Endpoint,
	)
	return &VolcEngine{client: client}
}

func (h *VolcEngine) CreateKnowledgeBase(c *gin.Context, req *types.KnowledgeBaseCreateReq) (*types.KnowledgeBaseCreateResp, error) {
	body := map[string]interface{}{
		"Name":        req.Name,
		"Description": req.Description,
	}

	result, err := h.client.Call("CreateKnowledgeBase", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("创建知识库失败: %v", err)
		return nil, exterr.Newf(-1, "创建知识库失败: %v", err)
	}

	kbID, _ := result["Result"].(map[string]interface{})["KnowledgeBaseId"].(string)
	if kbID == "" {
		return nil, exterr.Newf(-1, "创建知识库失败: 未返回知识库ID")
	}

	return &types.KnowledgeBaseCreateResp{KnowledgeBaseID: kbID}, nil
}

func (h *VolcEngine) ListKnowledgeBase(c *gin.Context, req *types.KnowledgeBaseListReq) (*types.KnowledgeBaseListResp, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	body := map[string]interface{}{
		"PageNum":  req.PageNum,
		"PageSize": req.PageSize,
	}

	result, err := h.client.Call("ListKnowledgeBase", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("查询知识库列表失败: %v", err)
		return nil, exterr.Newf(-1, "查询知识库列表失败: %v", err)
	}

	respData, ok := result["Result"].(map[string]interface{})
	if !ok {
		return nil, exterr.Newf(-1, "查询知识库列表失败: 响应格式错误")
	}

	total, _ := respData["Total"].(float64)
	listData, _ := respData["List"].([]interface{})

	list := make([]types.KnowledgeBaseInfo, 0, len(listData))
	for _, item := range listData {
		if m, ok := item.(map[string]interface{}); ok {
			info := types.KnowledgeBaseInfo{}
			if v, ok := m["KnowledgeBaseId"].(string); ok {
				info.ID = v
			}
			if v, ok := m["Name"].(string); ok {
				info.Name = v
			}
			if v, ok := m["Description"].(string); ok {
				info.Description = v
			}
			if v, ok := m["CreateTime"].(string); ok {
				info.CreatedAt = v
			}
			list = append(list, info)
		}
	}

	return &types.KnowledgeBaseListResp{
		Total: int64(total),
		List:  list,
	}, nil
}

func (h *VolcEngine) DeleteKnowledgeBase(c *gin.Context, req *types.KnowledgeBaseDeleteReq) (*types.KnowledgeBaseDeleteResp, error) {
	body := map[string]interface{}{
		"KnowledgeBaseId": req.KnowledgeBaseID,
	}

	_, err := h.client.Call("DeleteKnowledgeBase", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("删除知识库失败: %v", err)
		return nil, exterr.Newf(-1, "删除知识库失败: %v", err)
	}

	return &types.KnowledgeBaseDeleteResp{}, nil
}

func (h *VolcEngine) UploadDocument(c *gin.Context, req *types.DocumentUploadReq) (*types.DocumentUploadResp, error) {
	body := map[string]interface{}{
		"KnowledgeBaseId": req.KnowledgeBaseID,
		"FileName":        req.FileName,
		"FileUrl":         req.FileURL,
	}

	result, err := h.client.Call("UploadDocument", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("上传文档失败: %v", err)
		return nil, exterr.Newf(-1, "上传文档失败: %v", err)
	}

	docID, _ := result["Result"].(map[string]interface{})["DocumentId"].(string)
	if docID == "" {
		return nil, exterr.Newf(-1, "上传文档失败: 未返回文档ID")
	}

	return &types.DocumentUploadResp{DocumentID: docID}, nil
}

func (h *VolcEngine) ListDocument(c *gin.Context, req *types.DocumentListReq) (*types.DocumentListResp, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	body := map[string]interface{}{
		"KnowledgeBaseId": req.KnowledgeBaseID,
		"PageNum":         req.PageNum,
		"PageSize":        req.PageSize,
	}

	result, err := h.client.Call("ListDocument", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("查询文档列表失败: %v", err)
		return nil, exterr.Newf(-1, "查询文档列表失败: %v", err)
	}

	respData, ok := result["Result"].(map[string]interface{})
	if !ok {
		return nil, exterr.Newf(-1, "查询文档列表失败: 响应格式错误")
	}

	total, _ := respData["Total"].(float64)
	listData, _ := respData["List"].([]interface{})

	list := make([]types.DocumentInfo, 0, len(listData))
	for _, item := range listData {
		if m, ok := item.(map[string]interface{}); ok {
			info := types.DocumentInfo{}
			if v, ok := m["DocumentId"].(string); ok {
				info.ID = v
			}
			if v, ok := m["FileName"].(string); ok {
				info.FileName = v
			}
			if v, ok := m["FileUrl"].(string); ok {
				info.FileURL = v
			}
			if v, ok := m["Status"].(string); ok {
				info.Status = v
			}
			if v, ok := m["CreateTime"].(string); ok {
				info.CreatedAt = v
			}
			list = append(list, info)
		}
	}

	return &types.DocumentListResp{
		Total: int64(total),
		List:  list,
	}, nil
}

func (h *VolcEngine) DeleteDocument(c *gin.Context, req *types.DocumentDeleteReq) (*types.DocumentDeleteResp, error) {
	body := map[string]interface{}{
		"KnowledgeBaseId": req.KnowledgeBaseID,
		"DocumentId":      req.DocumentID,
	}

	_, err := h.client.Call("DeleteDocument", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("删除文档失败: %v", err)
		return nil, exterr.Newf(-1, "删除文档失败: %v", err)
	}

	return &types.DocumentDeleteResp{}, nil
}

func (h *VolcEngine) Chat(c *gin.Context, req *types.ChatReq) (*types.ChatResp, error) {
	if req.TopK <= 0 {
		req.TopK = 5
	}

	body := map[string]interface{}{
		"KnowledgeBaseId": req.KnowledgeBaseID,
		"Query":           req.Query,
		"TopK":            req.TopK,
	}

	result, err := h.client.Call("Chat", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("知识问答失败: %v", err)
		return nil, exterr.Newf(-1, "知识问答失败: %v", err)
	}

	respData, ok := result["Result"].(map[string]interface{})
	if !ok {
		return nil, exterr.Newf(-1, "知识问答失败: 响应格式错误")
	}

	answer, _ := respData["Answer"].(string)
	sourcesData, _ := respData["Sources"].([]interface{})

	sources := make([]types.Source, 0, len(sourcesData))
	for _, item := range sourcesData {
		if m, ok := item.(map[string]interface{}); ok {
			source := types.Source{}
			if v, ok := m["Content"].(string); ok {
				source.Content = v
			}
			if v, ok := m["Score"].(float64); ok {
				source.Score = v
			}
			sources = append(sources, source)
		}
	}

	return &types.ChatResp{
		Answer:  answer,
		Sources: sources,
	}, nil
}

func (h *VolcEngine) Search(c *gin.Context, req *types.SearchReq) (*types.SearchResp, error) {
	if req.TopK <= 0 {
		req.TopK = 10
	}

	body := map[string]interface{}{
		"KnowledgeBaseId": req.KnowledgeBaseID,
		"Query":           req.Query,
		"TopK":            req.TopK,
	}

	result, err := h.client.Call("Search", "2024-01-01", body)
	if err != nil {
		h.Log.Errorf("检索失败: %v", err)
		return nil, exterr.Newf(-1, "检索失败: %v", err)
	}

	respData, ok := result["Result"].(map[string]interface{})
	if !ok {
		return nil, exterr.Newf(-1, "检索失败: 响应格式错误")
	}

	resultsData, _ := respData["Results"].([]interface{})

	results := make([]types.SearchResult, 0, len(resultsData))
	for _, item := range resultsData {
		if m, ok := item.(map[string]interface{}); ok {
			res := types.SearchResult{}
			if v, ok := m["Content"].(string); ok {
				res.Content = v
			}
			if v, ok := m["Score"].(float64); ok {
				res.Score = v
			}
			results = append(results, res)
		}
	}

	return &types.SearchResp{Results: results}, nil
}
