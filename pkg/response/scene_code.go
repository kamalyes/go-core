/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-09 23:26:10
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-13 13:27:25
 * @FilePath: \go-core\pkg\response\scene_code.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import "sync"

// SceneCode 是自定义业务状态码类型
type SceneCode int

// ResponseContext 结构体
type ResponseContext struct {
	SceneCode SceneCode
	TrackId   string
}

// 自定义状态码
const (
	Success       = 200  // 成功
	BadRequest    = 400  // 错误请求
	Fail          = 500  // 失败
	ServerError   = 1000 // 服务器错误
	ValidateError = 1001 // 参数校验错误
	Deadline      = 1002 // 服务调用超时
	CreateError   = 1003 // 服务器写入失败
	FindError     = 1004 // 服务器查询失败
	WithoutServer = 1005 // 服务未启用
	AuthError     = 1006 // 权限错误
	DeleteError   = 1007 // 服务器删除失败
	EmptyFile     = 1008 // 文件为空
	RateLimit     = 1009 // 访问限流
	Unauthorized  = 1010 // JWT认证失败
	WithoutLogin  = 1011 // 用户未登录
	DisableAuth   = 1012 // 当前用户已被禁用
)

// sceneCodeMsgMap 用于存储状态码和消息的映射关系
var sceneCodeMsgMap = struct {
	sync.RWMutex
	mapping map[SceneCode]string
}{
	mapping: map[SceneCode]string{
		Success:       "Success",
		BadRequest:    "Bad Request",
		Fail:          "Fail",
		ServerError:   "Internal Server Error",
		ValidateError: "Validation Error",
		Deadline:      "Deadline Exceeded",
		CreateError:   "Failed to Create",
		FindError:     "Failed to Find",
		WithoutServer: "Service Unavailable",
		AuthError:     "Authorization Error",
		DeleteError:   "Failed to Delete",
		EmptyFile:     "Empty File",
		RateLimit:     "Rate Limit Exceeded",
		Unauthorized:  "Unauthorized",
		WithoutLogin:  "User Not Logged In",
		DisableAuth:   "User Authentication Disabled",
	},
}

// SetSceneCode 设置自定义状态码
func SetSceneCode(code SceneCode, msg string) {
	sceneCodeMsgMap.Lock()
	defer sceneCodeMsgMap.Unlock()
	sceneCodeMsgMap.mapping[code] = msg
}

// GetSceneCodeMsg 根据状态码获取对应的错误消息
func GetSceneCodeMsg(code SceneCode) string {
	sceneCodeMsgMap.RLock()
	defer sceneCodeMsgMap.RUnlock()
	return sceneCodeMsgMap.mapping[code]
}
