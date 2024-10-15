/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-15 20:13:06
 * @FilePath: \go-core\response\ginx.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kamalyes/go-core/global"
)

// GenGinResponse 生成 JSON 格式的响应
func GenGinResponse(ctx *gin.Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.merge()
	// 将 StatusCode 转换为 int 类型
	httpCodeInt := int(respOption.HttpCode)
	// 返回JSON格式的响应
	// 创建一个map来存储不包含HttpStatusCode和Language的字段
	cleanedResp := map[string]interface{}{
		"data":    respOption.Data,
		"code":    respOption.Code,
		"message": respOption.Message,
	}
	ctx.JSON(httpCodeInt, cleanedResp)
}

// GenGin400xResponse 生成 HTTP 400x 错误响应
func GenGin400xResponse(ctx *gin.Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = BadRequest
	respOption.HttpCode = StatusBadRequest
	GenGinResponse(ctx, respOption)
}

// GenGin500xResponse 生成 HTTP 500 错误响应
func GenGin500xResponse(ctx *gin.Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = Fail
	respOption.HttpCode = StatusInternalServerError
	GenGinResponse(ctx, respOption)
}

// GinValidatorError 处理字段校验异常
func GinValidatorError(ctx *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		GenGin400xResponse(ctx, &ResponseOption{
			Code: ValidateError,
			Data: removeTopStruct(errs.Translate(global.Trans)),
		})
		return
	}

	GenGin400xResponse(ctx, &ResponseOption{Message: err.Error()})
}

// removeTopStruct 定义一个去掉结构体名称前缀的自定义方法：
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
