package types

type KnowledgeBaseCreateReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type KnowledgeBaseCreateResp struct {
	KnowledgeBaseID string `json:"knowledgeBaseId"`
}

type KnowledgeBaseListReq struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}

type KnowledgeBaseListResp struct {
	Total int64               `json:"total"`
	List  []KnowledgeBaseInfo `json:"list"`
}

type KnowledgeBaseInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type KnowledgeBaseDeleteReq struct {
	KnowledgeBaseID string `json:"knowledgeBaseId" binding:"required"`
}

type KnowledgeBaseDeleteResp struct{}

type DocumentUploadReq struct {
	KnowledgeBaseID string `json:"knowledgeBaseId" binding:"required"`
	FileName        string `json:"fileName" binding:"required"`
	FileURL         string `json:"fileUrl" binding:"required"`
}

type DocumentUploadResp struct {
	DocumentID string `json:"documentId"`
}

type DocumentListReq struct {
	KnowledgeBaseID string `form:"knowledgeBaseId" binding:"required"`
	PageNum         int    `form:"pageNum"`
	PageSize        int    `form:"pageSize"`
}

type DocumentListResp struct {
	Total int64          `json:"total"`
	List  []DocumentInfo `json:"list"`
}

type DocumentInfo struct {
	ID        string `json:"id"`
	FileName  string `json:"fileName"`
	FileURL   string `json:"fileUrl"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type DocumentDeleteReq struct {
	KnowledgeBaseID string `json:"knowledgeBaseId" binding:"required"`
	DocumentID      string `json:"documentId" binding:"required"`
}

type DocumentDeleteResp struct{}

type ChatReq struct {
	KnowledgeBaseID string `json:"knowledgeBaseId" binding:"required"`
	Query           string `json:"query" binding:"required"`
	TopK            int    `json:"topK"`
}

type ChatResp struct {
	Answer  string   `json:"answer"`
	Sources []Source `json:"sources"`
}

type Source struct {
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type SearchReq struct {
	KnowledgeBaseID string `json:"knowledgeBaseId" binding:"required"`
	Query           string `json:"query" binding:"required"`
	TopK            int    `json:"topK"`
}

type SearchResp struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}
