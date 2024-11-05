# go-core

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/kamalyes/go-core)
[![license](https://img.shields.io/github/license/kamalyes/go-core)]()
[![download](https://img.shields.io/github/downloads/kamalyes/go-core/total)]()
[![release](https://img.shields.io/github/v/release/kamalyes/go-core)]()
[![commit](https://img.shields.io/github/last-commit/kamalyes/go-core)]()
[![issues](https://img.shields.io/github/issues/kamalyes/go-core)]()
[![pull](https://img.shields.io/github/issues-pr/kamalyes/go-core)]()
[![fork](https://img.shields.io/github/forks/kamalyes/go-core)]()
[![star](https://img.shields.io/github/stars/kamalyes/go-core)]()
[![go](https://img.shields.io/github/go-mod/go-version/kamalyes/go-core)]()
[![size](https://img.shields.io/github/repo-size/kamalyes/go-core)]()
[![contributors](https://img.shields.io/github/contributors/kamalyes/go-core)]()


### 介绍

go-core 是 go web 应用开发脚手架，从全局配置文件读取，zap日志组件始化，gorm数据库连接初始化，redis客户端初始化，http server启动等。最终实现简化流程、提高效率、统一规范。

### 安装

```bash
go get -u github.com/kamalyes/go-core
```

### 例子

默认的程序根目录下必须包含 resources 文件夹，且文件夹内必须有不同环境的开发文件至少一种
配置文件参考 <https://github.com/kamalyes/go-config> 库的resources目录下的配置文件

```shell
├── resources(项目整合配置文件示例)
│   └── dev_config.yaml  开发环境配置文件
│   └── fat_config.yaml  功能验收测试环境配置文件
│   └── pro_config.yaml  生产环境配置文件
│   └── uat_config.yaml  用户验收测试环境配置文件
```

```go
package main

import (
	"github.com/gin-gonic/gin"
	goconfig "github.com/kamalyes/go-config"
	"github.com/kamalyes/go-config/pkg/env"
	"github.com/kamalyes/go-core/pkg/database"
	"github.com/kamalyes/go-core/pkg/global"
	"github.com/kamalyes/go-core/pkg/minio"
	"github.com/kamalyes/go-core/pkg/mqtt"
	"github.com/kamalyes/go-core/pkg/redis"
	"github.com/kamalyes/go-core/pkg/srun"
	"github.com/kamalyes/go-core/pkg/zap"
)

func main() {
	// 获取全局配置,默认根据运行环境加载对应配置文件
	global.CONFIG = goconfig.GlobalConfig()

	// 初始化zap日志
	global.LOG = zap.Zap()

	// 初始化数据库连接
	global.DB = database.Gorm()

	// 初始化 redis 客户端
	global.REDIS = redis.Redis()

	// 初始化 minio
	global.MinIO = minio.Minio()

	// 初始化 mqtt
	global.MQTT = mqtt.DefaultMqtt("111111")

	// 获取配置文件原始内容,这样方便在程序中全局拿到自己定义的配置子项
	global.VP = global.CONFIG.Viper

	// 启动 http 服务
	r := gin.Default()
	// 健康监测
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, "ok")
	})
	// 启动服务
	srun.RunHttpServer(r)
}
```
