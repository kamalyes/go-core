/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-10 23:44:23
 * @FilePath: \go-core\pkg\response\ginx.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"github.com/gin-gonic/gin"
)

// ResponseOption 结构体定义
type ResponseOption struct {
	Data     interface{}  `json:"data"`
	Code     BusinessCode `json:"code"`
	HttpCode int
	Message  string `json:"message"`
}

// NewResponseOption 是 ResponseOption 的构造函数
func NewResponseOption(code BusinessCode) *ResponseOption {
	return &ResponseOption{
		Code:     code,
		Data:     nil,
		HttpCode: 200,
		Message:  GetBusinessCodeMsg(SUCCESS),
	}
}

// Sub 链式调用
func (o ResponseOption) Sub(ctx *gin.Context) {
	GenResponse(ctx, &o)
}

// GenResponse 函数用于返回响应
func GenResponse(ctx *gin.Context, respOption *ResponseOption) {
	respOpt := NewResponseOption(SUCCESS)
	if respOption != nil {
		respOpt.merge(respOption)
	}

	// 返回JSON格式的响应
	ctx.JSON(respOpt.HttpCode, ResponseOption{
		Code:    respOpt.Code,
		Data:    respOpt.Data,
		Message: respOpt.Message,
	})
}

// merge 方法用于合并非空字段
func (o *ResponseOption) merge(respOption *ResponseOption) {
	if respOption.Code != 0 {
		o.Code = respOption.Code
	}
	if respOption.Data != nil {
		o.Data = respOption.Data
	}
	if respOption.Message != "" {
		o.Message = respOption.Message
	}
	if respOption.HttpCode != 0 {
		o.HttpCode = respOption.HttpCode
	}
}
