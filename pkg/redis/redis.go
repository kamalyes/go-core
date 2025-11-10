/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-10 10:55:49
 * @FilePath: \go-core\pkg\redis\redis.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package redis

import (
	"context"

	"github.com/kamalyes/go-core/pkg/global"
	"github.com/redis/go-redis/v9"
)

// Redis 初始z化redis客户端
func Redis() *redis.Client {
	redisCfg := global.CONFIG.Redis
	if redisCfg.Addr == "" {
		return nil
	}
	db := 0
	if redisCfg.DB >= 0 && redisCfg.DB <= 15 {
		db = redisCfg.DB
	}
	client := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Addr,
		Password:     redisCfg.Password,
		DB:           db,
		MaxRetries:   redisCfg.MaxRetries,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
	})
	pong, err := client.Ping(context.TODO()).Result()
	if err != nil {
		global.LOGGER.ErrorKV("redis connect ping failed", "err", err)
		return nil
	} else {
		global.LOGGER.InfoKV("redis connect ping response", "pong", pong)
		return client
	}
}
