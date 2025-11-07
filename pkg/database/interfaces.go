/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 09:15:07
 * @FilePath: \go-core\pkg\database\interfaces.go
 * @Description: 数据库操作接口定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"database/sql"
	"gorm.io/gorm"
)

// QueryParam 定义查询参数接口
type QueryParam interface {
	Where(db *gorm.DB) *gorm.DB
}

// Handler 定义数据库操作接口
type Handler interface {
	DB() *gorm.DB
	Query(param QueryParam) *gorm.DB
	Close() error
	AutoMigrate(dst ...any) error
	Begin(opts ...*sql.TxOptions) Handler
	Commit() error
	Rollback() error
}