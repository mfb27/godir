package user

import (
	"fmt"

	"godir/internal/common/jwt"
	"godir/internal/common/svc"
	"godir/internal/handler/base"
	"godir/internal/model"
	"godir/internal/service"
	"godir/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type User struct {
	base.BaseHandler
	userService *service.User
}

func (h *User) New() base.BaseHandlerInterface {
	return &User{
		userService: service.NewUser(),
	}
}

func (h *User) Create(c *gin.Context, req *types.UserCreateReq) (*types.UserCreateResp, error) {
	if req.Source != 1 && req.Source != 2 {
		if h.Log != nil {
			h.Log.Info("param source is invaild")
		}
		return nil, fmt.Errorf("invalid source, must be 1 or 2")
	}
	user := model.User{
		Name:   req.Name,
		Source: req.Source,
	}
	err := h.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &types.UserCreateResp{
		UserId: user.ID,
	}, nil
}

// Profile 获取当前用户个人信息
func (h *User) Profile(c *gin.Context, req *types.UserProfileReq) (*types.UserProfileResp, error) {
	// 从上下文获取用户信息
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return nil, fmt.Errorf("无法获取用户信息")
	}

	claims, ok := userInfo.(jwt.Claims)
	if !ok {
		return nil, fmt.Errorf("用户信息格式错误")
	}

	// 获取用户详细信息及预签名头像URL
	user, avatarPresignedURL, err := h.userService.GetUserProfileWithPresignedAvatar(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 如果没有预签名URL，则使用原始URL
	finalAvatarURL := avatarPresignedURL
	if finalAvatarURL == "" {
		finalAvatarURL = user.Avatar
	}

	return &types.UserProfileResp{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   finalAvatarURL,
		Nickname: user.Nickname,
		Gender:   user.Gender,
	}, nil
}

// UpdateProfile 更新当前用户个人信息
func (h *User) UpdateProfile(c *gin.Context, req *types.UserProfileUpdateReq) (*types.UserProfileResp, error) {
	// 从上下文获取用户信息
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return nil, fmt.Errorf("无法获取用户信息")
	}

	claims, ok := userInfo.(jwt.Claims)
	if !ok {
		return nil, fmt.Errorf("用户信息格式错误")
	}

	// 更新用户信息
	err := h.userService.UpdateGodirUserProfile(claims.UserID, req.Avatar, req.Nickname, req.Gender)
	if err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	// 获取更新后的用户信息及预签名头像URL
	user, avatarPresignedURL, err := h.userService.GetUserProfileWithPresignedAvatar(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 如果没有预签名URL，则使用原始URL
	finalAvatarURL := avatarPresignedURL
	if finalAvatarURL == "" {
		finalAvatarURL = user.Avatar
	}

	return &types.UserProfileResp{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   finalAvatarURL,
		Nickname: user.Nickname,
		Gender:   user.Gender,
	}, nil
}

// UploadAvatar 上传用户头像
func (h *User) UploadAvatar(c *gin.Context, req *types.UploadAvatarReq) (*types.UploadAvatarResp, error) {
	// 从上下文获取用户信息
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return nil, fmt.Errorf("无法获取用户信息")
	}

	claims, ok := userInfo.(jwt.Claims)
	if !ok {
		return nil, fmt.Errorf("用户信息格式错误")
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("获取上传文件失败: %w", err)
	}

	// 生成文件名
	filename := fmt.Sprintf("avatars/%d_%d", claims.UserID, file.Size)

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 上传到MinIO
	_, err = svc.Minio().PutObject(
		c.Request.Context(),
		svc.Cfg().MinIO.Bucket,
		filename,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 更新数据库中的封面信息
	protocol := "http"
	if svc.Cfg().MinIO.UseSSL {
		protocol = "https"
	}
	// 生成预签名URL
	url := fmt.Sprintf("%s://%s/%s/%s", protocol, svc.Cfg().MinIO.Endpoint, svc.Cfg().MinIO.Bucket, filename)

	return &types.UploadAvatarResp{
		URL: url,
	}, nil
}
