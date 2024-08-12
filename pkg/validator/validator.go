/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 23:48:43
 * @FilePath: \go-core\pkg\validator\validator.go
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
	"github.com/kamalyes/go-core/pkg/response"
)

// HandleValidatorError 处理字段校验异常
func HandleValidatorError(ctx *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		response.GenResponse(ctx, &response.ResponseOption{
			Code: response.ValidateError,
			Data: removeTopStruct(errs.Translate(global.Trans)),
		})
		return
	}

	response.GenResponse(ctx, &response.ResponseOption{Message: err.Error()})
}

// removeTopStruct 定义一个去掉结构体名称前缀的自定义方法：
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
