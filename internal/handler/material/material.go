package material

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"godir/internal/common/jwt"
	"godir/internal/common/redis"
	"godir/internal/common/svc"
	"godir/internal/handler/base"
	"godir/internal/model"
	"godir/internal/types"

	"github.com/gin-gonic/gin"
)

type Material struct {
	base.BaseHandler
}

func (h *Material) New() base.BaseHandlerInterface {
	return new(Material)
}

// GetUploadToken 获取上传临时token
func (h *Material) GetUploadToken(c *gin.Context, req *types.MaterialUploadTokenReq) (*types.MaterialUploadTokenResp, error) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	cfg := svc.Cfg()

	// 生成文件key（使用用户ID和时间戳）
	ext := filepath.Ext(req.FileName)
	key := fmt.Sprintf("user/%d/%d%s", userIDUint, time.Now().UnixNano(), ext)

	// cred, err := miniox.TemporaryCredentials(cfg, []string{key})
	// if err != nil {
	// 	h.Log.Errorf("生成临时凭证失败: %v", err)
	// 	return nil, exterr.Fail
	// }

	// return &types.MaterialUploadTokenResp{
	// 	AccessKeyID:     cred.AccessKeyID,
	// 	SecretAccessKey: cred.SecretAccessKey,
	// 	SessionToken:    cred.SessionToken,
	// 	Bucket:          cfg.MinIO.Bucket,
	// 	Key:             key,
	// 	Endpoint:        miniox.BuildBaseUrl(cfg.MinIO.UseSSL, cfg.MinIO.Endpoint),
	// }, nil

	// 生成临时访问凭证
	// 注意：MinIO本身不支持STS，这里直接返回主凭证
	// 生产环境应该使用STS服务或自定义临时凭证生成
	protocol := "http"
	if cfg.MinIO.UseSSL {
		protocol = "https"
	}
	return &types.MaterialUploadTokenResp{
		AccessKeyID:     cfg.MinIO.AccessKeyID,
		SecretAccessKey: cfg.MinIO.SecretAccessKey,
		SessionToken:    "", // MinIO不支持session token
		Bucket:          cfg.MinIO.Bucket,
		Key:             key,
		Endpoint:        fmt.Sprintf("%s://%s", protocol, cfg.MinIO.Endpoint),
	}, nil
}

// Save 保存文件信息
func (h *Material) Save(c *gin.Context, req *types.MaterialSaveReq) (*types.MaterialSaveResp, error) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	cfg := svc.Cfg()
	protocol := "http"
	if cfg.MinIO.UseSSL {
		protocol = "https"
	}
	urlStr := fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.MinIO.Endpoint, req.Bucket, req.Key)

	// 创建文件记录
	material := model.GodirMaterial{
		UserID:      userIDUint,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ContentType: req.ContentType,
		Bucket:      req.Bucket,
		Key:         req.Key,
		URL:         urlStr,
	}

	if err := h.DB.Create(&material).Error; err != nil {
		return nil, fmt.Errorf("保存文件信息失败: %w", err)
	}

	// 索引到 Elasticsearch
	esClient := svc.ES()
	if esClient != nil {
		// 构造索引文档
		doc := map[string]interface{}{
			"id":           material.ID,
			"user_id":      material.UserID,
			"file_name":    material.FileName,
			"file_size":    material.FileSize,
			"content_type": material.ContentType,
			"bucket":       material.Bucket,
			"key":          material.Key,
			"url":          material.URL,
			"cover_key":    material.CoverKey,
			"cover_url":    material.CoverURL,
			"created_at":   material.CreatedAt,
		}

		b, err := json.Marshal(doc)
		if err != nil {
			h.Log.Warnf("ES 索引序列化失败: %v", err)
		} else {
			resp, err := esClient.Index(
				"godir_material",
				bytes.NewReader(b),
				esClient.Index.WithDocumentID(fmt.Sprintf("%d", material.ID)),
			)
			if err != nil {
				h.Log.Warnf("ES索引失败: %v", err)
			} else if resp != nil {
				defer resp.Body.Close()
			}
		}
	}

	err := redis.PushThumbnailTask(&redis.ThumbnailTask{
		MaterialID:  material.ID,
		Bucket:      req.Bucket,
		Key:         req.Key,
		ContentType: req.ContentType,
	})
	if err != nil {
		h.Log.Warnf("推送缩略图任务失败: %v", err)
	}

	// // generate thumbnail for image/video using ffmpeg, upload to MinIO and save cover info
	// ct := req.ContentType
	// if ct == "" {
	// 	ct = material.ContentType
	// }
	// if strings.HasPrefix(ct, "image/") || strings.HasPrefix(ct, "video/") {
	// 	minioClient := svc.Minio()
	// 	if minioClient != nil {
	// 		ctx := context.Background()

	// 		// download object to temp file
	// 		obj, err := minioClient.GetObject(ctx, req.Bucket, req.Key, minioLib.GetObjectOptions{})
	// 		if err == nil {
	// 			tmpFile, err := os.CreateTemp(os.TempDir(), "material-*")
	// 			if err == nil {
	// 				_, _ = io.Copy(tmpFile, obj)
	// 				_ = tmpFile.Close()
	// 				defer os.Remove(tmpFile.Name())

	// 				// make sure ffmpeg exists
	// 				if _, err := exec.LookPath("ffmpeg"); err == nil {
	// 					thumbPath := tmpFile.Name() + ".thumb.jpg"

	// 					var cmd *exec.Cmd
	// 					if strings.HasPrefix(ct, "image/") {
	// 						cmd = exec.Command("ffmpeg", "-y", "-i", tmpFile.Name(), "-vf", "scale=640:-1", "-vframes", "1", "-q:v", "2", thumbPath)
	// 					} else {
	// 						// video
	// 						cmd = exec.Command("ffmpeg", "-y", "-i", tmpFile.Name(), "-ss", "00:00:01", "-vframes", "1", "-q:v", "2", thumbPath)
	// 					}

	// 					if out, err := cmd.CombinedOutput(); err == nil {
	// 						// upload thumbnail
	// 						tf, err := os.Open(thumbPath)
	// 						if err == nil {
	// 							fi, _ := tf.Stat()
	// 							thumbKey := req.Key + ".thumb.jpg"
	// 							_, err = minioClient.PutObject(ctx, req.Bucket, thumbKey, tf, fi.Size(), minioLib.PutObjectOptions{ContentType: "image/jpeg"})
	// 							_ = tf.Close()
	// 							if err == nil {
	// 								protocol := "http"
	// 								if svc.Cfg().MinIO.UseSSL {
	// 									protocol = "https"
	// 								}
	// 								coverURL := fmt.Sprintf("%s://%s/%s/%s", protocol, svc.Cfg().MinIO.Endpoint, req.Bucket, thumbKey)
	// 								material.CoverKey = thumbKey
	// 								material.CoverURL = coverURL
	// 								_ = h.DB.Model(&material).Updates(map[string]interface{}{"cover_key": material.CoverKey, "cover_url": material.CoverURL}).Error
	// 							}
	// 							_ = os.Remove(thumbPath)
	// 						}
	// 					} else {
	// 						// ffmpeg failed; keep silent (could log out)
	// 						_ = out
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return &types.MaterialSaveResp{MaterialID: material.ID}, nil
}

// List 获取文件列表
func (h *Material) List(c *gin.Context, req *types.MaterialListReq) (*types.MaterialListResp, error) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	var materials []model.GodirMaterial
	if err := h.DB.Where("user_id = ?", userIDUint).
		Order("created_at DESC").
		Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("查询文件列表失败: %w", err)
	}

	materialList := make([]types.MaterialInfo, 0, len(materials))
	for _, m := range materials {
		// 生成预签名URL
		downloadURL := ""
		previewURL := ""
		coverPreviewURL := ""

		minioClient := svc.Minio()
		if minioClient != nil {
			// 生成7天有效的预签名URL用于下载
			reqParams := make(map[string][]string)
			reqParams["response-content-disposition"] = []string{fmt.Sprintf("attachment; filename=\"%s\"", m.FileName)}

			presigned, err := minioClient.PresignedGetObject(
				c.Request.Context(),
				m.Bucket,
				m.Key,
				time.Hour*24*7, // 7天有效期
				reqParams)
			if err == nil {
				downloadURL = presigned.String()
			}

			// 生成预览URL（内联显示）
			inlineParams := make(map[string][]string)
			inlineParams["response-content-disposition"] = []string{"inline"}

			getURL, err := minioClient.PresignedGetObject(
				c.Request.Context(),
				m.Bucket,
				m.Key,
				time.Hour*24*7, // 7天有效期
				inlineParams)
			if err == nil {
				previewURL = getURL.String()
			}

			// cover preview URL (inline) for thumbnail display
			if m.CoverKey != "" {
				coverParams := make(map[string][]string)
				coverParams["response-content-disposition"] = []string{"inline"}
				coverURL, err := minioClient.PresignedGetObject(
					c.Request.Context(),
					m.Bucket,
					m.CoverKey,
					time.Hour*24*7,
					coverParams,
				)
				if err == nil {
					coverPreviewURL = coverURL.String()
				}
			}
		}

		materialList = append(materialList, types.MaterialInfo{
			ID:              m.ID,
			FileName:        m.FileName,
			FileSize:        m.FileSize,
			ContentType:     m.ContentType,
			URL:             m.URL,
			CoverURL:        m.CoverURL,
			CoverPreviewURL: coverPreviewURL, // 前端可用此字段展示封面
			DownloadURL:     downloadURL,
			PreviewURL:      previewURL,
			CreatedAt:       m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.MaterialListResp{
		Materials: materialList,
	}, nil
}

// Search 根据素材名称搜索（优先使用ES，未配置则回退到DB模糊查询）
func (h *Material) Search(c *gin.Context, req *types.MaterialSearchReq) (*types.MaterialSearchResp, error) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	q := strings.TrimSpace(req.Q)
	if q == "" {
		return &types.MaterialSearchResp{Materials: []types.MaterialInfo{}}, nil
	}

	var results []types.MaterialInfo

	esClient := svc.ES()
	if esClient != nil {
		// 使用 ES 搜索，按文件名前缀匹配并过滤当前用户
		body := map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []interface{}{
						map[string]interface{}{"match_phrase_prefix": map[string]interface{}{"file_name": map[string]interface{}{"query": q}}},
						map[string]interface{}{"term": map[string]interface{}{"user_id": userIDUint}},
					},
				},
			},
			"size": 50,
		}

		b, _ := json.Marshal(body)
		resp, err := esClient.Search(
			esClient.Search.WithIndex("godir_material"),
			esClient.Search.WithBody(bytes.NewReader(b)),
		)
		if err == nil && resp != nil {
			defer resp.Body.Close()
			var r map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&r); err == nil {
				if hitsObj, ok := r["hits"].(map[string]any); ok {
					if hitsArr, ok := hitsObj["hits"].([]any); ok {
						for _, hi := range hitsArr {
							if hit, ok := hi.(map[string]any); ok {
								if src, ok := hit["_source"].(map[string]any); ok {
									var mi types.MaterialInfo
									// 尝试从_source中映射字段
									if v, ok := src["id"].(float64); ok {
										mi.ID = uint(v)
									}
									if v, ok := src["file_name"].(string); ok {
										mi.FileName = v
									}
									if v, ok := src["file_size"].(float64); ok {
										mi.FileSize = int64(v)
									}
									if v, ok := src["content_type"].(string); ok {
										mi.ContentType = v
									}
									if v, ok := src["url"].(string); ok {
										mi.URL = v
									}
									if v, ok := src["cover_url"].(string); ok {
										mi.CoverURL = v
									}
									if v, ok := src["created_at"].(string); ok {
										mi.CreatedAt = v
									}

									// 补充 preview/download URL（如 MinIO 可用，则使用预签名）
									if svc.Minio() != nil {
										// 尝试从 URL 提取 bucket/key
										if u, err := url.Parse(mi.URL); err == nil {
											cfg := svc.Cfg()
											if u.Host == cfg.MinIO.Endpoint && u.Path != "" {
												parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
												if len(parts) >= 2 {
													bucket := parts[0]
													objectKey := strings.Join(parts[1:], "/")
													inlineParams := make(map[string][]string)
													inlineParams["response-content-disposition"] = []string{"inline"}
													if presigned, err := svc.Minio().PresignedGetObject(c.Request.Context(), bucket, objectKey, time.Hour*24, inlineParams); err == nil {
														mi.PreviewURL = presigned.String()
													}
													reqParams := make(map[string][]string)
													reqParams["response-content-disposition"] = []string{fmt.Sprintf("attachment; filename=\"%s\"", mi.FileName)}
													if presigned2, err := svc.Minio().PresignedGetObject(c.Request.Context(), bucket, objectKey, time.Hour*24, reqParams); err == nil {
														mi.DownloadURL = presigned2.String()
													}
												}
											}
										}
									}

									results = append(results, mi)
								}
							}
						}
					}
				}
			}
		}
	}

	// // 如果 ES 未启用或未返回结果，则回退到数据库模糊查询
	// if len(results) == 0 {
	// 	var materials []model.GodirMaterial
	// 	if err := h.DB.Where("user_id = ? AND file_name LIKE ?", userIDUint, "%"+q+"%").Order("created_at DESC").Limit(50).Find(&materials).Error; err == nil {
	// 		for _, m := range materials {
	// 			mi := types.MaterialInfo{
	// 				ID:          m.ID,
	// 				FileName:    m.FileName,
	// 				FileSize:    m.FileSize,
	// 				ContentType: m.ContentType,
	// 				URL:         m.URL,
	// 				CoverURL:    m.CoverURL,
	// 				CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
	// 			}

	// 			// 生成预签名 URL（若 MinIO 可用）
	// 			if svc.Minio() != nil && (strings.HasPrefix(mi.ContentType, "image/") || strings.HasPrefix(mi.ContentType, "video/")) {
	// 				inlineParams := make(map[string][]string)
	// 				inlineParams["response-content-disposition"] = []string{"inline"}
	// 				if presigned, err := svc.Minio().PresignedGetObject(c.Request.Context(), m.Bucket, m.Key, time.Hour*24, inlineParams); err == nil {
	// 					mi.PreviewURL = presigned.String()
	// 				}
	// 			}

	// 			results = append(results, mi)
	// 		}
	// 	}
	// }

	return &types.MaterialSearchResp{Materials: results}, nil
}

// Publish publishes a material with a description
func (h *Material) Publish(c *gin.Context, req *types.MaterialPublishReq) (*types.MaterialPublishResp, error) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	// 检查素材是否存在且属于当前用户
	var material model.GodirMaterial
	result := h.DB.Where("id = ? AND user_id = ?", req.MaterialID, userIDUint).First(&material)
	if result.Error != nil {
		return nil, fmt.Errorf("素材不存在或无权限操作: %w", result.Error)
	}

	// 创建发布记录
	published := &model.GodirPublishedMaterial{
		UserID:      userIDUint,
		MaterialID:  req.MaterialID,
		Description: req.Description,
	}

	result = h.DB.Create(published)
	if result.Error != nil {
		return nil, fmt.Errorf("发布素材失败: %w", result.Error)
	}

	resp := &types.MaterialPublishResp{
		ID: published.ID,
	}

	return resp, nil
}

// LikePublish 点赞某条发布
func (h *Material) LikePublish(c *gin.Context, req *types.PublishLikeReq) (*types.PublishLikeResp, error) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	// 检查发布记录是否存在
	var published model.GodirPublishedMaterial
	if err := h.DB.Where("id = ?", req.PublishID).First(&published).Error; err != nil {
		return nil, fmt.Errorf("发布不存在: %w", err)
	}

	// 尝试查找是否已点赞
	var like model.GodirPublishedLike
	if err := h.DB.Where("published_id = ? AND user_id = ?", req.PublishID, userIDUint).First(&like).Error; err == nil {
		// 已存在，直接返回计数
	} else {
		// 创建点赞记录（若并发可能出现重复唯一键错误，忽略该错误）
		_ = h.DB.Create(&model.GodirPublishedLike{UserID: userIDUint, PublishedID: req.PublishID}).Error
	}

	// 统计点赞数
	var count int64
	h.DB.Model(&model.GodirPublishedLike{}).Where("published_id = ?", req.PublishID).Count(&count)

	return &types.PublishLikeResp{LikesCount: int(count), Liked: true}, nil
}

// UnlikePublish 取消点赞
func (h *Material) UnlikePublish(c *gin.Context, req *types.PublishLikeReq) (*types.PublishLikeResp, error) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	// 删除点赞记录（若不存在也无妨）
	_ = h.DB.Where("published_id = ? AND user_id = ?", req.PublishID, userIDUint).Delete(&model.GodirPublishedLike{}).Error

	// 统计点赞数
	var count int64
	h.DB.Model(&model.GodirPublishedLike{}).Where("published_id = ?", req.PublishID).Count(&count)

	return &types.PublishLikeResp{LikesCount: int(count), Liked: false}, nil
}

// BatchDelete 批量删除文件（同时支持单个和多个文件删除）
func (h *Material) BatchDelete(c *gin.Context, req *types.MaterialBatchDeleteReq) (*types.MaterialBatchDeleteResp, error) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	if len(req.Ids) == 0 {
		return nil, fmt.Errorf("请选择要删除的文件")
	}

	// 查询要删除的文件
	var materials []model.GodirMaterial
	if err := h.DB.Where("id IN (?) AND user_id = ?", req.Ids, userIDUint).Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}

	// 检查是否所有请求的文件都找到了
	foundIds := make(map[uint]bool)
	for _, material := range materials {
		foundIds[material.ID] = true
	}

	for _, id := range req.Ids {
		if !foundIds[id] {
			return nil, fmt.Errorf("文件ID %d 不存在或无权限删除", id)
		}
	}

	// 从数据库中删除记录（注意：不删除MinIO中的实际文件）
	if err := h.DB.Where("id IN (?)", req.Ids).Delete(&model.GodirMaterial{}).Error; err != nil {
		return nil, fmt.Errorf("删除文件记录失败: %w", err)
	}

	return &types.MaterialBatchDeleteResp{}, nil
}

// ListPublished 获取发布列表
func (h *Material) ListPublished(c *gin.Context, req *types.PublishListReq) (*types.PublishListResp, error) {
	var publishedMaterials []model.GodirPublishedMaterial

	// 查询所有发布记录，按创建时间倒序排列
	if err := h.DB.Order("created_at DESC").Find(&publishedMaterials).Error; err != nil {
		return nil, fmt.Errorf("查询发布列表失败: %w", err)
	}

	// 构造响应数据
	// 获取请求上下文中的当前用户（若存在），用于判断是否已点赞
	var currentUserID uint
	if uidI, ok := c.Get("userId"); ok {
		if uid, ok2 := uidI.(uint); ok2 {
			currentUserID = uid
		}
	}

	// 如果没有从中间件获得 userId，则尝试从 Authorization 头解析 JWT（允许前端在 public 接口中携带 token）
	if currentUserID == 0 {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if claims, err := jwt.ParseToken(tokenStr); err == nil && claims != nil {
				currentUserID = claims.UserID
			}
		}
	}

	list := make([]types.PublishInfo, 0, len(publishedMaterials))
	for _, published := range publishedMaterials {
		info := types.PublishInfo{
			ID:          published.ID,
			UserID:      published.UserID,
			MaterialID:  published.MaterialID,
			Description: published.Description,
			CreatedAt:   published.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// 获取关联的用户信息
		var user model.GodirUser
		if err := h.DB.Where("id = ?", published.UserID).First(&user).Error; err == nil {
			avatarUrl := user.Avatar

			// 如果头像是存储在MinIO中的资源，则生成预签名URL
			if avatarUrl != "" && svc.Minio() != nil {
				// 解析头像URL，检查是否是MinIO中的资源
				if u, err := url.Parse(avatarUrl); err == nil {
					// 检查是否是MinIO的endpoint
					cfg := svc.Cfg()
					if u.Host == cfg.MinIO.Endpoint && u.Path != "" {
						// 从路径中提取bucket和object key
						pathParts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
						if len(pathParts) >= 2 {
							bucket := pathParts[0]
							objectKey := strings.Join(pathParts[1:], "/")

							// 生成预签名URL
							presignedURL, err := svc.Minio().PresignedGetObject(
								c,
								bucket,
								objectKey,
								time.Hour*24, // 24小时有效期
								nil)
							if err == nil {
								avatarUrl = presignedURL.String()
							}
						}
					}
				}
			}

			info.User = types.UserInfo{
				ID:       user.ID,
				Username: user.Username,
				Nickname: user.Nickname,
				Avatar:   avatarUrl,
			}
		}

		// 获取关联的素材信息
		var material model.GodirMaterial
		if err := h.DB.Where("id = ?", published.MaterialID).First(&material).Error; err == nil {
			// 生成预签名URL用于前端展示
			var previewUrl string
			var coverPreviewUrl string

			// 为图片和视频文件生成预览URL
			if strings.HasPrefix(material.ContentType, "image/") || strings.HasPrefix(material.ContentType, "video/") {
				// 生成预览URL（内联显示）
				inlineParams := make(url.Values)
				inlineParams.Set("response-content-disposition", "inline")

				presignedURL, err := svc.Minio().PresignedGetObject(
					c,
					material.Bucket,
					material.Key,
					time.Hour*24, // 24小时有效期
					inlineParams)
				if err == nil {
					previewUrl = presignedURL.String()
				}

				// 如果有封面，也生成封面的预览URL
				if material.CoverKey != "" {
					coverParams := make(url.Values)
					coverParams.Set("response-content-disposition", "inline")
					coverPresignedURL, err := svc.Minio().PresignedGetObject(
						c,
						material.Bucket,
						material.CoverKey,
						time.Hour*24,
						coverParams,
					)
					if err == nil {
						coverPreviewUrl = coverPresignedURL.String()
					}
				}
			}

			materialInfo := types.MaterialInfo{
				ID:              material.ID,
				FileName:        material.FileName,
				FileSize:        material.FileSize,
				ContentType:     material.ContentType,
				URL:             material.URL,
				CoverURL:        material.CoverURL,
				PreviewURL:      previewUrl,
				CoverPreviewURL: coverPreviewUrl,
				CreatedAt:       material.CreatedAt.Format("2006-01-02 15:04:05"),
			}

			info.Material = materialInfo
		}

		// 统计点赞数并判断当前用户是否已点赞
		var likesCount int64
		h.DB.Model(&model.GodirPublishedLike{}).Where("published_id = ?", published.ID).Count(&likesCount)
		info.LikesCount = int(likesCount)

		if currentUserID != 0 {
			var like model.GodirPublishedLike
			if err := h.DB.Where("published_id = ? AND user_id = ?", published.ID, currentUserID).First(&like).Error; err == nil {
				info.Liked = true
			}
		}

		list = append(list, info)
	}

	resp := &types.PublishListResp{
		List: list,
	}

	return resp, nil
}

// UpdateMaterialName 修改素材文件名
func (h *Material) UpdateMaterialName(c *gin.Context, req *types.MaterialUpdateNameReq) (*types.MaterialUpdateNameResp, error) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		return nil, fmt.Errorf("未登录")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return nil, fmt.Errorf("用户ID格式错误")
	}

	// 校验新文件名
	newName := strings.TrimSpace(req.NewName)
	if newName == "" {
		return nil, fmt.Errorf("文件名不能为空")
	}

	// 查询素材是否存在且属于当前用户
	var material model.GodirMaterial
	result := h.DB.Where("id = ? AND user_id = ?", req.MaterialID, userIDUint).First(&material)
	if result.Error != nil {
		return nil, fmt.Errorf("素材不存在或无权限操作: %w", result.Error)
	}

	// 更新文件名
	if err := h.DB.Model(&material).Update("file_name", newName).Error; err != nil {
		return nil, fmt.Errorf("更新文件名失败: %w", err)
	}

	// 如果启用了 ES，也更新索引
	esClient := svc.ES()
	if esClient != nil {
		doc := map[string]interface{}{
			"file_name": newName,
		}

		// ES Update API expects a body like { "doc": { ... } }
		updateBody := map[string]interface{}{
			"doc": doc,
		}

		b, err := json.Marshal(updateBody)
		if err != nil {
			h.Log.Warnf("ES 文档序列化失败: %v", err)
		} else {
			resp, err := esClient.Update(
				"godir_material",
				fmt.Sprintf("%d", material.ID),
				bytes.NewReader(b),
			)
			if err != nil {
				h.Log.Warnf("ES 更新失败: %v", err)
			} else if resp != nil {
				defer resp.Body.Close()
				if resp.IsError() {
					h.Log.Warnf("ES 响应错误: %v", resp.String())
				} else {
					h.Log.Infof("ES 更新成功: %v", resp.String())
				}
			}
		}
	}

	return &types.MaterialUpdateNameResp{
		MaterialID: material.ID,
		FileName:   newName,
	}, nil
}
