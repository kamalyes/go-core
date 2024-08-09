/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:16:39
 * @FilePath: \go-core\zap\syncer .go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package zap

import (
	"os"

	"github.com/kamalyes/go-core/global"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// WriteSyncer 利用lumberjack库做日志分割
func WriteSyncer(file string) zapcore.WriteSyncer {
	// 在进行切割之前，日志文件的最大大小（以MB为单位）
	maxSize := 10
	if global.CONFIG.Zap.MaxSize > 10 && global.CONFIG.Zap.MaxSize < 500 {
		maxSize = global.CONFIG.Zap.MaxSize
	}
	// 保留旧文件的最大个数
	maxBackups := 100
	if global.CONFIG.Zap.MaxBackups > 100 {
		maxBackups = global.CONFIG.Zap.MaxBackups
	}
	// 保留旧文件的最大天数
	maxAge := 30
	if global.CONFIG.Zap.MaxAge > 30 {
		maxAge = global.CONFIG.Zap.MaxBackups
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   global.CONFIG.Zap.Compress,
	}
	// 是否需要在控制台输出日志
	if global.CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
