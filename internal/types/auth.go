package types

// 登录接口
type (
	AuthLoginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	AuthLoginResp struct {
		Token    string `json:"token"`
		UserID   uint   `json:"userId"`
		Username string `json:"username"`
	}
)

// 登出接口
type (
	AuthLogoutReq struct {
	}
	
	AuthLogoutResp struct {
		Message string `json:"message"`
	}
)

// 注册接口
type (
	AuthRegisterReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	
	AuthRegisterResp struct {
		UserID   uint   `json:"userId"`
		Username string `json:"username"`
		Message  string `json:"message"`
	}
)