package types

type AiAppListReq struct{}

type AiAppInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	AppID string `json:"appId"`
	Desc  string `json:"desc"`
	Icon  string `json:"icon"`
	Cover string `json:"cover"`
}

type AiAppListResp struct {
	Apps []AiAppInfo `json:"apps"`
}
