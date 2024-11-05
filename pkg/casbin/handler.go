/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:12:15
 * @FilePath: \go-core\pkg\casbin\handler.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package casbin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/pkg/global"
	"github.com/kamalyes/go-core/pkg/jwt"
)

var casbinAdmi = global.GPerFix + "casbin_admi"

// SetCasbinAdmi 设置casbin_admi常量的值
func SetCasbinAdmi(value string) {
	casbinAdmi = value
}

// GetCasbinAdmi 获取casbin_admi常量的值
func GetCasbinAdmi() string {
	return casbinAdmi
}

// CasbinHandler Casbin权限认证
func CasbinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.Contains(path, "swagger") {
			ctx.Next()
			return
		}
		if strings.Contains(path, "login") || strings.Contains(path, "health") || strings.Contains(path, "captcha") {
			ctx.Next()
			return
		}
		// 从ctx中获取claims
		claims, _ := jwt.GetClaims(ctx)
		user := claims.UserId
		permission := ctx.Request.URL.Path
		method := ctx.Request.Method

		if claims.UserType != GetCasbinAdmi() {
			ok := CasbinServiceApp.PermissionVerify(user, permission, method)
			if !ok {
				ctx.JSON(http.StatusMethodNotAllowed, gin.H{
					"code":    -1,
					"message": "用户已经通过身份验证，但请求的接口:(" + permission + ")不在您的权限之内！",
				})
				ctx.Abort()
				return
			}
		}
		ctx.Next()
		return
	}
}
