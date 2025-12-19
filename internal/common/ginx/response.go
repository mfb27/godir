package ginx

import "godir/internal/common/exterr"

type Response struct {
	Code int64  `json:"code"` // 0成功 其他失败
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) Response {
	return Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
}

// Fail 导出fail函数供其他包使用
func Fail(err error) Response {
	return fail(err)
}

func fail(err error) Response {
	return Response{
		Code: exterr.Code(err),
		Msg:  exterr.Msg(err),
	}
}
