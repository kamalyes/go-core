/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 13:55:26
 * @FilePath: \go-core\internal\response\response.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"net/http"
	"sync"

	"github.com/kamalyes/go-core/internal/code"

	"github.com/gin-gonic/gin"
)

var responseMapKey = "km_http_code"

// SetResponseMapKey 设置 responseMapKey 变量的值
func SetResponseMapKey(value string) {
	responseMapKey = value
}

// GetResponseMapKey 获取 responseMapKey 变量的值
func GetResponseMapKey() string {
	return responseMapKey
}

// SetHTTPCode 设置默认的 httpCode
func SetHTTPCode(code int) {
	responseMap.Store(GetResponseMapKey(), code)
}

var responseMap sync.Map

// Response 统一 json 结构体
type Response struct {

	/** 状态码 */
	Code code.StatusCode `json:"code"`

	/** 内容体 */
	Data interface{} `json:"data"`

	/** 消息 */
	Message string `json:"message"`
}

// Result gin 统一返回
func Result(ctx *gin.Context, code code.StatusCode, data interface{}, message string) {
	var httpCodeX int
	httpCodeInterface, ok := responseMap.Load(GetResponseMapKey())
	if ok {
		httpCodeX = httpCodeInterface.(int)
		// 删除对应键，重置responseMap
		responseMap.Delete(GetResponseMapKey())
	}
	if code > 0 {
		httpCodeX = http.StatusOK
	}
	ctx.JSON(httpCodeX, Response{
		code,
		data,
		message,
	})
}

// Ok 成功
func Ok(ctx *gin.Context) {
	Result(ctx, code.SUCCESS, map[string]interface{}{}, code.GetErrMsg(code.SUCCESS))
}

// OkMsg 带message消息的成功
func OkMsg(ctx *gin.Context, message string) {
	Result(ctx, code.SUCCESS, map[string]interface{}{}, message)
}

// OkData 带数据的成功
func OkData(ctx *gin.Context, data interface{}) {
	Result(ctx, code.SUCCESS, data, code.GetErrMsg(code.SUCCESS))
}

// OkDataMsg 带数据和返回消息的成功
func OkDataMsg(ctx *gin.Context, data interface{}, message string) {
	Result(ctx, code.SUCCESS, data, message)
}

// Fail 失败
func Fail(ctx *gin.Context) {
	Result(ctx, code.FAIL, map[string]interface{}{}, code.GetErrMsg(code.FAIL))
}

// FailMsg 带message消息的失败
func FailMsg(ctx *gin.Context, message string) {
	Result(ctx, code.FAIL, map[string]interface{}{}, message)
}

// Fail 失败
func FailData(ctx *gin.Context, statusCode code.StatusCode, data interface{}, message ...string) {
	messageX := code.GetErrMsg(statusCode)
	if len(message) > 0 {
		messageX = message[0]
	}
	Result(ctx, statusCode, data, messageX)
}

// FailDataMsg 带数据和返回消息的失败
func FailDataMsg(ctx *gin.Context, data interface{}, message string) {
	Result(ctx, code.FAIL, data, message)
}
