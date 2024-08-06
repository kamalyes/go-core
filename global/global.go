/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-06 10:58:48
 * @FilePath: \go-core\global\global.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package global

import (
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/casbin/casbin/v2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	goconfig "github.com/kamalyes/go-config"
	"github.com/kamalyes/go-config/env"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

// 单节点应用常用全局变量
var (

	// ENV 设置环境
	ENV env.Environment

	// DB 数据库
	DB *gorm.DB

	// REDIS 默认客户端
	REDIS *redis.Client

	// MQTT 客户端
	MQTT *mqtt.Client

	// CONFIG 全局系统配置
	CONFIG *goconfig.Config

	// VP 通过 viper 读取的yaml配置文件
	VP *viper.Viper

	// LOG 全局日志
	LOG *zap.Logger

	// CSBEF casbin实施者
	CSBEF casbin.IEnforcer

	// 雪花ID节点
	Node *snowflake.Node

	// MinIO客户端
	MinIO *minio.Client
)

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

func InitZapLogger() {
	osHostname := GetHostname()
	initialFields := map[string]interface{}{
		"hostname": osHostname,
		"pid":      os.Getpid(),
	}
	// 创建配置
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel), // 设置日志级别
		Encoding:         "json",                              // 设置日志格式
		OutputPaths:      []string{"stdout"},                  // 输出位置，这里输出到标准输出
		ErrorOutputPaths: []string{"stderr"},                  // 错误输出位置
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		InitialFields: initialFields,
	}

	// 根据配置构建Logger
	logger, err := cfg.Build()
	if err != nil {
		// 配置错误，处理异常
		panic(err)
	}

	LOG = logger
}
