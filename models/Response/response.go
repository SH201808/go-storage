package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Params gin.H  `json:"data"`
}

func basicResponse(code int, msg string, params gin.H) Response {
	response := Response{
		Code:   code,
		Msg:    msg,
		Params: params,
	}

	return response
}

func Success(data ...interface{}) Response {
	code := 0
	message := "成功"
	params := gin.H{}

	for _, v := range data {
		switch v.(type) {
		case int:
			code = v.(int)
		case string:
			message = v.(string)
		case gin.H:
			params = v.(gin.H)
		}
	}
	return basicResponse(code, message, params)
}

func Err(data ...interface{}) Response {
	code := -1
	message := "错误"
	params := gin.H{}

	for _, value := range data {
		switch value.(type) {
		case int:
			code = value.(int)
		case string:
			message = value.(string)
		case gin.H:
			params = value.(gin.H)
		}
	}

	return basicResponse(code, message, params)
}
