/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-09 23:26:10
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-10 22:39:46
 * @FilePath: \go-core\pkg\response\common.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import "sync"

// BusinessCode 是自定义业务状态码类型
type BusinessCode int

// ResponseContext 结构体
type ResponseContext struct {
	BusinessCode BusinessCode
	TrackId      string
}

// 自定义状态码
const (
	SUCCESS       = 200
	FAIL          = 500
	ServerError   = 1000
	ValidateError = 1001
	Deadline      = 1002
	CreateError   = 1003
	FindError     = 1004
	WithoutServer = 1005
	AuthError     = 1006
	DeleteError   = 1007
	EmptyFile     = 1008
	RateLimit     = 1009
	Unauthorized  = 10010
	WithoutLogin  = 10011
	DisableAuth   = 10012
)

// StatusCodeMsgMap 用于存储状态码和消息的映射关系
var BusinessCodeMsgMap = struct {
	sync.RWMutex
	mapping map[BusinessCode]string
}{
	mapping: map[BusinessCode]string{
		SUCCESS:       "成功",
		FAIL:          "失败",
		ServerError:   "服务器错误",
		ValidateError: "参数校验错误",
		Deadline:      "服务调用超时",
		CreateError:   "服务器写入失败",
		FindError:     "服务器查询失败",
		WithoutServer: "服务未启用",
		AuthError:     "权限错误",
		DeleteError:   "服务器删除失败",
		EmptyFile:     "文件为空",
		RateLimit:     "访问限流",
		Unauthorized:  "JWT认证失败",
		WithoutLogin:  "用户未登录",
		DisableAuth:   "当前用户已被禁用",
	},
}

// SetBusinessCode 设置自定义状态码
func SetBusinessCode(code BusinessCode, msg string) {
	BusinessCodeMsgMap.Lock()
	defer BusinessCodeMsgMap.Unlock()

	BusinessCodeMsgMap.mapping[code] = msg
}

// GetBusinessCodeMsg 根据状态码获取对应的错误消息
func GetBusinessCodeMsg(code BusinessCode) string {
	BusinessCodeMsgMap.RLock()
	defer BusinessCodeMsgMap.RUnlock()

	return BusinessCodeMsgMap.mapping[code]
}
