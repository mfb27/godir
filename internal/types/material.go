package types

// 上传令牌接口
type (
	MaterialUploadTokenReq struct {
		FileName    string `json:"fileName" binding:"required"`
		FileSize    int64  `json:"fileSize" binding:"required"`
		ContentType string `json:"contentType"`
	}

	MaterialUploadTokenResp struct {
		AccessKeyID     string `json:"accessKeyId"`
		SecretAccessKey string `json:"secretAccessKey"`
		SessionToken    string `json:"sessionToken"`
		Bucket          string `json:"bucket"`
		Key             string `json:"key"`
		Endpoint        string `json:"endpoint"`
	}
)

// 保存材料接口
type (
	MaterialSaveReq struct {
		FileName    string `json:"fileName" binding:"required"`
		FileSize    int64  `json:"fileSize" binding:"required"`
		ContentType string `json:"contentType"`
		Bucket      string `json:"bucket" binding:"required"`
		Key         string `json:"key" binding:"required"`
		URL         string `json:"url" binding:"required"`
	}

	MaterialSaveResp struct {
		MaterialID uint `json:"materialId"`
	}
)

// 材料列表接口
type (
	MaterialListReq  struct{}
	MaterialListResp struct {
		Materials []MaterialInfo `json:"materials"`
	}

	// 搜索请求
	MaterialSearchReq struct {
		Q string `form:"q" binding:"required" json:"q"`
	}
	MaterialSearchResp struct {
		Materials []MaterialInfo `json:"materials"`
	}

	MaterialInfo struct {
		ID              uint   `json:"id"`
		FileName        string `json:"fileName"`
		FileSize        int64  `json:"fileSize"`
		ContentType     string `json:"contentType"`
		URL             string `json:"url"`
		CoverURL        string `json:"coverUrl"`
		CoverPreviewURL string `json:"coverPreviewUrl"` // 预签名封面URL（用于页面展示）
		DownloadURL     string `json:"downloadUrl"`     // 预签名下载URL
		PreviewURL      string `json:"previewUrl"`      // 预签名预览URL
		CreatedAt       string `json:"createdAt"`
	}
)

// 批量删除接口
type (
	MaterialBatchDeleteReq struct {
		Ids []uint `json:"ids" binding:"required"`
	}

	MaterialBatchDeleteResp struct {
		Message string `json:"message"`
	}
)

// Publish material request
type MaterialPublishReq struct {
	MaterialID  uint   `json:"materialId"`
	Description string `json:"description"`
}

// Publish material response
type MaterialPublishResp struct {
	ID uint `json:"id"`
}

// PublishInfo represents published material info
type PublishInfo struct {
	ID          uint         `json:"id"`
	UserID      uint         `json:"userId"`
	MaterialID  uint         `json:"materialId"`
	Description string       `json:"description"`
	CreatedAt   string       `json:"createdAt"`
	User        UserInfo     `json:"user,omitempty"` // 用户信息结构体
	Material    MaterialInfo `json:"material,omitempty"`
	LikesCount  int          `json:"likesCount"`
	Liked       bool         `json:"liked,omitempty"` // 当前用户是否已点赞（若有用户上下文）
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// PublishListResp represents publish list response
type (
	PublishListReq  struct{}
	PublishListResp struct {
		List []PublishInfo `json:"list"`
	}
)

// Publish like/unlike
type (
	PublishLikeReq struct {
		PublishID uint `json:"publishId" binding:"required"`
	}

	PublishLikeResp struct {
		LikesCount int  `json:"likesCount"`
		Liked      bool `json:"liked"`
	}
)

// 修改素材文件名接口
type (
	MaterialUpdateNameReq struct {
		MaterialID uint   `json:"materialId" binding:"required"`
		NewName    string `json:"newName" binding:"required"`
	}

	MaterialUpdateNameResp struct {
		MaterialID uint   `json:"materialId"`
		FileName   string `json:"fileName"`
	}
)
