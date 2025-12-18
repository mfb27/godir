package auth

import (
	"fmt"

	"godir/internal/common/jwt"
	"godir/internal/handler/base"
	"godir/internal/model"
	"godir/internal/types"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	base.BaseHandler
}

func (h *Auth) New() base.BaseHandlerInterface {
	return new(Auth)
}

// Login 用户登录
func (h *Auth) Login(c *gin.Context, req *types.AuthLoginReq) (*types.AuthLoginResp, error) {
	// 查询用户
	var user model.GodirUser
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	return &types.AuthLoginResp{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

// Register 用户注册
func (h *Auth) Register(c *gin.Context, req *types.AuthRegisterReq) (*types.AuthRegisterResp, error) {
	// 检查用户名是否已存在
	var existingUser model.GodirUser
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := model.GodirUser{
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &types.AuthRegisterResp{
		UserID:   user.ID,
		Username: user.Username,
		Message:  "注册成功",
	}, nil
}

// Logout 用户退出
func (h *Auth) Logout(c *gin.Context, req *types.AuthLogoutReq) (*types.AuthLogoutResp, error) {
	// JWT是无状态的，退出主要是客户端删除token
	// 如果需要服务端控制，可以实现token黑名单机制
	return &types.AuthLogoutResp{
		Message: "退出成功",
	}, nil
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword 验证密码
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
