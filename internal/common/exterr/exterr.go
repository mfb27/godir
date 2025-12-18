package exterr

import (
	"fmt"
)

var Fail = New(-1, "fail")

type exterr struct {
	code int64
	msg  string
}

func (e *exterr) Error() string {
	return e.msg
}

func (e *exterr) Code() int64 {
	return e.code
}

func New(code int64, msg string) error {
	return &exterr{code: code, msg: msg}
}

func Newf(code int64, format string, args ...any) error {
	return New(code, fmt.Sprintf(format, args...))
}

// func Is(err error, code int64) bool {
// 	return errors.Is(err, &exterr{code: code})
// }

// func As(err error, code int64) *exterr {
// 	var e *exterr
// 	if errors.As(err, &e) && e.code == code {
// 		return e
// 	}
// 	return nil
// }

func Code(err error) int64 {
	if v, ok := err.(*exterr); ok {
		return v.code
	}
	return -1
}

func Msg(err error) string {
	if v, ok := err.(*exterr); ok {
		return v.msg
	}
	return "fail"
}
