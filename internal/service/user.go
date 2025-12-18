package service

import (
	"context"
	"net/url"
	"strings"
	"time"

	"godir/internal/common/svc"
	"godir/internal/model"
)

type User struct{}

// NewUser 创建用户服务实例
func NewUser() *User {
	return &User{}
}

// Create 创建用户
func (s *User) Create(name string) (*model.User, error) {
	user := &model.User{
		Name: name,
	}

	if err := svc.DB().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID 根据ID获取用户
func (s *User) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := svc.DB().First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetGodirUserByID 根据ID获取GodirUser用户
func (s *User) GetGodirUserByID(id uint) (*model.GodirUser, error) {
	var user model.GodirUser
	if err := svc.DB().First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateGodirUserProfile 更新GodirUser用户个人信息
func (s *User) UpdateGodirUserProfile(id uint, avatar, nickname string, gender int) error {
	updates := map[string]interface{}{}
	
	if avatar != "" {
		updates["avatar"] = avatar
	}
	
	if nickname != "" {
		updates["nickname"] = nickname
	}
	
	// 性别始终更新（可以设置为0）
	updates["gender"] = gender
	
	if err := svc.DB().Model(&model.GodirUser{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	
	return nil
}

// GetUserProfileWithPresignedAvatar 获取用户信息并生成头像预签名URL
func (s *User) GetUserProfileWithPresignedAvatar(id uint) (*model.GodirUser, string, error) {
	user, err := s.GetGodirUserByID(id)
	if err != nil {
		return nil, "", err
	}

	// 如果没有头像，直接返回
	if user.Avatar == "" {
		return user, "", nil
	}

	// 解析头像URL获取bucket和object key
	u, err := url.Parse(user.Avatar)
	if err != nil {
		// 如果解析失败，返回原始URL
		return user, user.Avatar, nil
	}

	// 从路径中提取bucket和key
	// 假设URL格式为 http://minio:9000/bucket/key
	path := u.Path
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// 分割路径获取bucket和key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		// 如果格式不匹配，返回原始URL
		return user, user.Avatar, nil
	}

	bucket := parts[0]
	key := parts[1]

	// 生成预签名URL
	presignedURL, err := svc.Minio().PresignedGetObject(
		context.Background(),
		bucket,
		key,
		time.Hour*24*7, // 7天有效期
		nil,
	)
	if err != nil {
		// 如果生成失败，返回原始URL
		return user, user.Avatar, nil
	}

	return user, presignedURL.String(), nil
}