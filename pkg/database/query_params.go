/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 09:58:19
 * @FilePath: \go-core\pkg\database\query_params.go
 * @Description: 查询参数实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"gorm.io/gorm"
)

// SimpleQueryParam 简单查询参数实现
type SimpleQueryParam struct {
	whereClause string
	args        []interface{}
}

// NewSimpleQueryParam 创建简单查询参数
func NewSimpleQueryParam(where string, args ...interface{}) QueryParam {
	return &SimpleQueryParam{
		whereClause: where,
		args:        args,
	}
}

// Where 实现 QueryParam 接口
func (q *SimpleQueryParam) Where(db *gorm.DB) *gorm.DB {
	if q.whereClause == "" {
		return db
	}
	return db.Where(q.whereClause, q.args...)
}

// PageQueryParam 分页查询参数
type PageQueryParam struct {
	whereClause string
	args        []interface{}
	limit       int
	offset      int
	orderBy     string
}

// NewPageQueryParam 创建分页查询参数
func NewPageQueryParam(where string, args []interface{}, limit, offset int, orderBy string) QueryParam {
	return &PageQueryParam{
		whereClause: where,
		args:        args,
		limit:       limit,
		offset:      offset,
		orderBy:     orderBy,
	}
}

// Where 实现 QueryParam 接口
func (q *PageQueryParam) Where(db *gorm.DB) *gorm.DB {
	if q.whereClause != "" {
		db = db.Where(q.whereClause, q.args...)
	}
	if q.orderBy != "" {
		db = db.Order(q.orderBy)
	}
	if q.limit > 0 {
		db = db.Limit(q.limit)
	}
	if q.offset > 0 {
		db = db.Offset(q.offset)
	}
	return db
}
