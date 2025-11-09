//go:build windows
// +build windows

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 18:32:24
 * @FilePath: \go-core\pkg\srun\windows.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package srun

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/pkg/global"
)

// RunHttpServer Windows环境下启动服务
func RunHttpServer(r *gin.Engine) {
	address := global.CONFIG.Server.Endpoint
	s := initServer(address, r)
	// 保证文本能够顺序输出
	time.Sleep(20 * time.Microsecond)
	global.LOGGER.InfoKV("server run success on", "address", address)
	err := s.ListenAndServe()
	if err != nil {
		global.LOGGER.Error(err.Error())
	}
}

// 初始化服务
func initServer(address string, router *gin.Engine) server {
	s := &http.Server{
		Addr:           address,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}
	return s
}
