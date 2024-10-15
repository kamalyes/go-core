package response

import (
	"github.com/labstack/echo/v4"
)

// GenEchoResponse 生成 JSON 格式的响应
func GenEchoResponse(ctx echo.Context, respOption *ResponseOption) error {
	return SendJSONResponse(ctx, respOption)
}

// GenEcho400xResponse 生成 HTTP 400x 错误响应
func GenEcho400xResponse(ctx echo.Context, respOption *ResponseOption) error {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = BadRequest
	respOption.HttpCode = StatusBadRequest
	return GenEchoResponse(ctx, respOption)
}

// GenEcho500xResponse 生成 HTTP 500 错误响应
func GenEcho500xResponse(ctx echo.Context, respOption *ResponseOption) error {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = Fail
	respOption.HttpCode = StatusInternalServerError
	return GenEchoResponse(ctx, respOption)
}
