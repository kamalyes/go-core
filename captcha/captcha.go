/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 13:56:35
 * @FilePath: \go-core\captcha\captcha.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package captcha

import (
	"context"
	"time"

	"github.com/kamalyes/go-core/global"
	"go.uber.org/zap"
)

var (
	expirationTime = time.Second * 180
	perFixKey      = global.GPerFix + "captcha"
)

type RedisStore struct {
	Expiration time.Duration
	PrefixKey  string
	Context    context.Context
}

type RedisError string

// RedisStoreInterface 是 RedisStore 的接口
type RedisStoreInterface interface {
	Set(id string, value string) error
	Get(key string, clear bool) (string, error)
	Verify(id, answer string, clear bool) bool
	UseWithCtx(ctx context.Context) *RedisStore
}

// SetExpirationTime 设置过期时间
func SetExpirationTime(duration time.Duration) {
	expirationTime = duration
}

// GetExpirationTime 获取过期时间
func GetExpirationTime() time.Duration {
	return expirationTime
}

// SetPerFixKey 设置前缀键
func SetPerFixKey(key string) {
	perFixKey = key
}

// GetPerFixKey 获取前缀键
func GetPerFixKey() string {
	return perFixKey
}

// NewDefaultRedisStore 初始化默认的验证码 Redis 共享存储器
func NewDefaultRedisStore() RedisStoreInterface {
	return &RedisStore{
		Expiration: GetExpirationTime(),
		PrefixKey:  GetPerFixKey(),
		Context:    context.Background(),
	}
}

// UseWithCtx 设置 RedisStore 的上下文
func (rs *RedisStore) UseWithCtx(ctx context.Context) *RedisStore {
	rs.Context = ctx
	return rs
}

// logRedisStoreError 记录 Redis 操作错误并返回错误
func logRedisStoreError(action string, err error) error {
	global.LOG.Error(action, zap.Error(err))
	return err
}

// Set 在 Redis 中设置键值对
func (rs *RedisStore) Set(key string, value string) error {
	err := global.REDIS.Set(rs.Context, rs.PrefixKey+key, value, rs.Expiration).Err()
	if err != nil {
		return logRedisStoreError("RedisStoreSetError!", err)
	}
	return nil
}

// Get 从 Redis 中获取键值对，并可选择清除
func (rs *RedisStore) Get(key string, clear bool) (string, error) {
	val, err := global.REDIS.Get(rs.Context, key).Result()
	if err != nil {
		return val, logRedisStoreError("RedisStoreGetError!", err)
	}
	if clear {
		err := global.REDIS.Del(rs.Context, key).Err()
		if err != nil {
			return val, logRedisStoreError("RedisStoreClearError!", err)
		}
	}
	return val, nil
}

// Verify 验证验证码答案是否正确
func (rs *RedisStore) Verify(key, answer string, clear bool) bool {
	v, err := rs.Get(rs.PrefixKey+key, clear)
	if err != nil {
		return false
	}
	return v == answer
}
