/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-08 15:35:19
 * @FilePath: \go-core\internal\validator\validator.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kamalyes/go-core/global"
	"github.com/kamalyes/go-core/internal/code"
	"github.com/kamalyes/go-core/internal/response"
)

// HandleValidatorError 处理字段校验异常
func HandleValidatorError(ctx *gin.Context, err error) {
	//如何返回错误信息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		response.FailMsg(ctx, err.Error())
	}
	response.FailData(ctx, code.ValidateError, removeTopStruct(errs.Translate(global.Trans)))
}

// removeTopStruct 定义一个去掉结构体名称前缀的自定义方法：
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
