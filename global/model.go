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
)

type DistributedId int64
type TTime time.Time

type Model struct {
	ID         DistributedId `json:"id,omitempty"            gorm:"column:id;primary_key;"`
	CreateTime TTime         `json:"createTime,omitempty"    gorm:"column:create_time;comment:创建时间;"`
	UpdateTime TTime         `json:"updateTime,omitempty"    gorm:"column:update_time;comment:更新时间;"`
}

// CreateId
/**
 *  @Description: 创建一个分布式ID（雪花ID）
 *  @return DistributedId
 */
func CreateId() DistributedId {
	id := Node.Generate()
	return DistributedId(id.Int64())
}

// CreateTime
/**
 *  @Description: 创建一个时间戳
 *  @return Time
 */
func CreateTime() TTime {
	t := time.Now()
	tTime := TTime(t)
	return tTime
}
