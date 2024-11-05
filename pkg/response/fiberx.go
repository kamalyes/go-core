/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-15 20:14:35
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-15 20:15:07
 * @FilePath: \go-core\pkg\response\fiberx.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"github.com/gofiber/fiber/v2"
)

// GenFiberResponse 生成 JSON 格式的响应
func GenFiberResponse(ctx *fiber.Ctx, respOption *ResponseOption) error {
	return SendJSONResponse(ctx, respOption)
}

// GenFiber400xResponse 生成 HTTP 400x 错误响应
func GenFiber400xResponse(ctx *fiber.Ctx, respOption *ResponseOption) error {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = BadRequest
	respOption.HttpCode = StatusBadRequest
	return GenFiberResponse(ctx, respOption)
}

// GenFiber500xResponse 生成 HTTP 500 错误响应
func GenFiber500xResponse(ctx *fiber.Ctx, respOption *ResponseOption) error {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = Fail
	respOption.HttpCode = StatusInternalServerError
	return GenFiberResponse(ctx, respOption)
}
