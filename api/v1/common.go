package v1

import (
	"encoding/json"
	"gin_mall/serizlizer"
)

func ErrorResponse(err error) serizlizer.Response {
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serizlizer.Response{
			Status: 400,
			Msg:    "JSON类型不匹配",
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: 400,
		Msg:    "请求参数错误",
		Error:  err.Error(),
	}
}
