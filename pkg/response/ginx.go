/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 19:26:08
 * @FilePath: \go-core\pkg\response\ginx.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"github.com/gin-gonic/gin"
)

type ResponseOption struct {
	Data     interface{}  `json:"data"`
	Code     BusinessCode `json:"code"`
	HttpCode int
	Message  string `json:"message"`
}

func (o *ResponseOption) Sub(ctx *gin.Context) {
	GenResponse(ctx, o)
}

func (o *ResponseOption) merge() *ResponseOption {
	if o.Code == 0 {
		o.Code = SUCCESS
	}
	if o.HttpCode == 0 {
		o.HttpCode = SUCCESS
	}
	if o.Message == "" {
		o.Message = GetBusinessCodeMsg(SUCCESS)
	}
	return o
}

func GenResponse(ctx *gin.Context, respOption *ResponseOption) {
	respOption.merge()
	// 返回JSON格式的响应
	ctx.JSON(respOption.HttpCode, ResponseOption{
		Code:    respOption.Code,
		Data:    respOption.Data,
		Message: respOption.Message,
	})
}
