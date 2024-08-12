package main

import (
	"github.com/gin-gonic/gin"
	goconfig "github.com/kamalyes/go-config"
	"github.com/kamalyes/go-config/env"
	"github.com/kamalyes/go-core/database"
	"github.com/kamalyes/go-core/global"
	"github.com/kamalyes/go-core/minio"
	"github.com/kamalyes/go-core/mqtt"
	"github.com/kamalyes/go-core/redis"
	"github.com/kamalyes/go-core/srun"
	"github.com/kamalyes/go-core/zap"
)

func main() {

	// 获取程序运行环境，默认会读取 resources/active.yaml 文件中配置的运行环境
	global.ENV = env.Active()

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
