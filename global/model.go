/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-08 15:34:14
 * @FilePath: \go-core\global\model.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package global

import (
	"time"

	"github.com/kamalyes/go-core/internal/dtype"
)

type Model struct {
	ID         dtype.DistributedId `json:"id,omitempty"            gorm:"column:id;primary_key;"`
	CreateTime dtype.Time          `json:"createTime,omitempty"    gorm:"column:create_time;comment:创建时间;"`
	UpdateTime *dtype.Time         `json:"updateTime,omitempty"    gorm:"column:update_time;comment:更新时间;"`
}

// CreateId
/**
 *  @Description: 创建一个分布式ID（雪花ID）
 *  @return DistributedId
 */
func CreateId() dtype.DistributedId {
	id := Node.Generate()
	return dtype.DistributedId(id.Int64())
}

// CreateTime
/**
 *  @Description: 创建一个时间戳
 *  @return Time
 */
func CreateTime() dtype.Time {
	t := time.Now()
	tTime := dtype.Time(t)
	return tTime
}
