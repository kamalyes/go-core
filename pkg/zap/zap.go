/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-03 20:24:51
 * @FilePath: \go-core\pkg\zap\zap.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package zap

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kamalyes/go-core/pkg/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	lowercaseLevelEncoder      = "LowercaseLevelEncoder"      // 小写级别编码器
	lowercaseColorLevelEncoder = "LowercaseColorLevelEncoder" // 小写级别带颜色编码器
	capitalLevelEncoder        = "CapitalLevelEncoder"        // 大写级别编码器
	capitalColorLevelEncoder   = "CapitalColorLevelEncoder"   // 大写级别带颜色编码器
)

func Zap() *zap.Logger {
	if ok, err := logDirectorExists(global.CONFIG.Zap.Director); !ok {
		fmt.Printf("创建目录：%v\n", global.CONFIG.Zap.Director)
		if err != nil {
			fmt.Println("目录检查失败:", err)
		}
		_ = os.Mkdir(global.CONFIG.Zap.Director, os.ModePerm)
	}

	cores := []zapcore.Core{
		getEncoderCore("server_debug.log", zap.DebugLevel),
		getEncoderCore("server_info.log", zap.InfoLevel),
		getEncoderCore("server_warn.log", zap.WarnLevel),
		getEncoderCore("server_error.log", zap.ErrorLevel),
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())

	if global.CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return logger
}

func getEncoderConfig() zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: global.CONFIG.Zap.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    customTimeEncoder,
		EncodeCaller:  zapcore.FullCallerEncoder,
	}
	switch {
	case global.CONFIG.Zap.EncodeLevel == lowercaseLevelEncoder:
		// 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case global.CONFIG.Zap.EncodeLevel == lowercaseColorLevelEncoder:
		// 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case global.CONFIG.Zap.EncodeLevel == capitalLevelEncoder:
		// 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case global.CONFIG.Zap.EncodeLevel == capitalColorLevelEncoder:
		// 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

func getEncoder() zapcore.Encoder {
	if global.CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderCore(fileName string, level zapcore.Level) zapcore.Core {
	writer := WriteSyncer(global.CONFIG.Zap.Director + fileName)
	return zapcore.NewCore(getEncoder(), writer, level)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(global.CONFIG.Zap.Prefix + "2006/01/02 - 15:04:05.000"))
}

func logDirectorExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.New("无法获取目录信息")
	}

	if !fi.IsDir() {
		return false, errors.New("目录已存在同名文件")
	}

	return true, nil
}
