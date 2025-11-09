//go:build !windows
// +build !windows

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 19:41:52
 * @FilePath: \go-core\pkg\srun\others.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package srun

import (
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/pkg/global"
)

// RunHttpServer Linux，unix等环境下启动服务
func RunHttpServer(r *gin.Engine) {
	address := global.CONFIG.Server.Endpoint
	s := initServer(address, r)
	// 保证文本顺序输出
	time.Sleep(20 * time.Microsecond)
	global.LOGGER.InfoKV("server run success on", "address", address)
	err := s.ListenAndServe()
	if err != nil {
		global.LOGGER.Error(err.Error())
	}
}

// 初始化服务
func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
