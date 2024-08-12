/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-13 16:09:59
 * @FilePath: \go-core\pkg\response\ginx.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/pkg/httpx"
)

// ResponseOption 是用于构建返回响应的结构体
type ResponseOption struct {
	Data       interface{}
	Code       SceneCode
	StatusCode httpx.StatusCode
	Message    string
}

// convertToSceneCode 辅助函数用于将输入值转换为 SceneCode 类型
func convertToSceneCode(val interface{}) SceneCode {
	code, ok := val.(SceneCode)
	if !ok {
		// 如果类型断言失败，可以执行适当的错误处理逻辑
		// 这里简单地返回一个默认值
		return SceneCode(Success)
	}
	return SceneCode(code)
}

// convertToHttpStatusCode 辅助函数用于将输入值转换为 StatusCode 类型
func convertToHttpStatusCode(val interface{}) httpx.StatusCode {
	statusCode, ok := val.(httpx.StatusCode)
	if !ok {
		return httpx.StatusCode(httpx.StatusOK)
	}
	return httpx.StatusCode(statusCode)
}

// NewResponseOption 用于创建 ResponseOption 实例
func NewResponseOption(data interface{}, options ...interface{}) *ResponseOption {
	response := &ResponseOption{
		Data: data,
	}

	for _, option := range options {
		switch opt := option.(type) {
		case SceneCode:
			response.Code = opt
		case httpx.StatusCode:
			response.StatusCode = opt
		case string:
			response.Message = opt

		}
	}

	return response
}

// Sub 用于在给定的上下文中生成响应
func (o *ResponseOption) Sub(ctx *gin.Context) {
	GenResponse(ctx, o)
}

// merge 用于处理 ResponseOption 实例的属性值
func (o *ResponseOption) merge() *ResponseOption {
	// 将 o.Code 的值根据条件进行转换
	o.Code = convertToSceneCode(ternary(o.Code == 0, Success, o.Code))

	// 根据条件将 o.StatusCode 的值进行转换
	o.StatusCode = convertToHttpStatusCode(ternary(o.StatusCode == 0, httpx.StatusOK, o.StatusCode))

	// 根据条件设置消息内容
	if o.Message == "" {
		o.Message = GetSceneCodeMsg(o.Code)
	}
	if o.Message == "" {
		o.Message = httpx.GetStatusCodeText(o.StatusCode)
	}

	return o
}

// ternary 函数实现三元运算
func ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// GenResponse 生成 JSON 格式的响应
func GenResponse(ctx *gin.Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.merge()
	// 将 StatusCode 转换为 int 类型
	statusCodeInt := int(respOption.StatusCode)
	// 返回JSON格式的响应
	// 创建一个map来存储不包含HttpStatusCode和Language的字段
	cleanedResp := map[string]interface{}{
		"data":    respOption.Data,
		"code":    respOption.Code,
		"message": respOption.Message,
	}
	ctx.JSON(statusCodeInt, cleanedResp)
}

// Gen400xResponse 生成 HTTP 400x 错误响应
func Gen400xResponse(ctx *gin.Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = BadRequest
	respOption.StatusCode = httpx.StatusBadRequest
	GenResponse(ctx, respOption)
}

// Gen500xResponse 生成 HTTP 500 错误响应
func Gen500xResponse(ctx *gin.Context, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = Fail
	respOption.StatusCode = httpx.StatusInternalServerError
	GenResponse(ctx, respOption)
}
