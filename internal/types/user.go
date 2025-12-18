package types

type (
	UserCreateReq struct {
		Name   string `json:"name"`
		Source int64  `json:"source"`
	}
	UserCreateResp struct {
		UserId uint `json:"userId"`
	}

	// UserProfileUpdateReq 更新用户个人信息请求
	UserProfileUpdateReq struct {
		Avatar   string `json:"avatar,omitempty"`
		Nickname string `json:"nickname,omitempty"`
		Gender   int    `json:"gender,omitempty"`
	}
)
type (
	UserProfileReq struct{}
	// UserProfileResp 用户个人信息响应
	UserProfileResp struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
		Gender   int    `json:"gender"`
	}
)

type (
	UploadAvatarReq  struct{}
	UploadAvatarResp struct {
		URL string `json:"url"`
	}
)
