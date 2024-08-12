/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 09:18:56
 * @FilePath: \go-core\minio\minio.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package minio

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kamalyes/go-core/global"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	once    sync.Once
	timeOut = 10
)

// Minio 初始z化minio客户端
func Minio() *minio.Client {
	minioCfg := global.CONFIG.Minio
	if minioCfg.Url == "" {
		return nil
	}
	endpoint := fmt.Sprintf("%s:%d", minioCfg.Url, minioCfg.Port)
	// 创建MinIO客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(minioCfg.AccessKey, minioCfg.SecretKey, ""),
	})
	if err != nil {
		global.LOG.Error("minio new client failed, err:", zap.Any("err", err))
	}
	// 检查服务状态
	_, err = client.HealthCheck(3)
	if err != nil {
		global.LOG.Error("minio connect ping failed, err:", zap.Any("err", err))
		return nil
	} else {
		return client
	}
}

// watch
/**
 *  @Description: 单独的监控协程
 *  @param ctx 上下文
 *  @param status 状态标记
 *  @param describe 异常描述
 *  @param timeOut 超时时间
 */
func watch(ctx context.Context, status *bool, describe string, timeOut int) {
	for {
		select {
		case <-ctx.Done():
			if !*status {
				panic(fmt.Sprintf("%s: 连接超时(%ds)", describe, timeOut))
			}
			return
		}
	}
}

// initBucket
/**
 *  @Description: 初始化存储桶
 *  @param bucketName
 *  @param location
 *  @return err
 */
func initBucket(bucketName, location string) (err error) {
	init := false
	// 创建一个子节点的context,自动超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	go watch(ctx, &init, "初始化对象存储minio失败", timeOut)
	if location == "" {
		location = "us-east-8"
	}
	exists, errBucketExists := global.MinIO.BucketExists(ctx, bucketName)
	if errBucketExists == nil && exists {
		log.Printf("存储桶%s已存在", bucketName)
	} else {
		err = global.MinIO.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			err = errors.New("创建存储桶" + bucketName + "失败")
		} else {
			log.Printf("创建存储桶%s成功", bucketName)
		}
	}
	init = true
	log.Println("初始化Minio成功")
	return
}

// UpLoad
/**
 *  @Description: 上传
 *  @param objectName
 *  @param filePath
 *  @param contentType application/zip
 */
func UpLoad(objectName, filePath, contentType string) (err error) {
	init := false
	// 创建一个子节点的context,自动超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	go watch(ctx, &init, "上传数据到对象存储minio失败", timeOut)
	bucketName := global.VP.GetString("MinIO.bucketName")
	once.Do(func() {
		err = initBucket(bucketName, "")
		if err != nil {
			panic(err)
		}
	})
	_, err = global.MinIO.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	init = true
	return
}

// DownLoad
/**
 *  @Description: 下载
 *  @param objectName
 *  @param filePath
 *  @param contentType
 */
func DownLoad(objectName, filePath string) (err error) {
	init := false
	// 创建一个子节点的context,自动超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	go watch(ctx, &init, "从对象存储minio下载数据失败", timeOut)

	bucketName := global.VP.GetString("MinIO.bucketName")
	err = global.MinIO.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})

	init = true
	return
}

// Delete
/**
 *  @Description: 删除文件
 *  @param objectName
 *  @return err
 */
func Delete(objectName string) (err error) {
	init := false
	// 创建一个子节点的context,自动超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	go watch(ctx, &init, "从对象存储minio删除数据失败", timeOut)

	bucketName := global.VP.GetString("MinIO.bucketName")
	//删除一个文件
	err = global.MinIO.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{GovernanceBypass: true})

	init = true
	return
}

// Exist
/**
 *  @Description: 文件是否存在
 *  @param objectName
 *  @param filePath
 *  @param contentType
 */
func Exist(objectName string) (exist bool) {
	init := false
	// 创建一个子节点的context,自动超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	go watch(ctx, &init, "判断对象存储minio存在数据失败", timeOut)

	bucketName := global.VP.GetString("MinIO.bucketName")
	_, err := global.MinIO.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	init = true
	return true
}
