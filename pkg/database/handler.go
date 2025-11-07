/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 09:15:07
 * @FilePath: \go-core\pkg\database\handler.go
 * @Description: 数据库处理器实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"database/sql"

	"gorm.io/gorm"
)

// DatabaseHandler 实现 Handler 接口
type DatabaseHandler struct {
	db *gorm.DB
}

// NewHandler 创建新的数据库处理器
func NewHandler(db *gorm.DB) Handler {
	return &DatabaseHandler{db: db}
}

// DB implements Handler
func (d *DatabaseHandler) DB() *gorm.DB {
	return d.db
}

// Query implements Handler
func (d *DatabaseHandler) Query(param QueryParam) *gorm.DB {
	return param.Where(d.db)
}

// Close implements Handler
func (d *DatabaseHandler) Close() error {
	sqlDB, _ := d.db.DB()
	return sqlDB.Close()
}

// AutoMigrate implements Handler
func (d *DatabaseHandler) AutoMigrate(dst ...any) error {
	return d.db.AutoMigrate(dst...)
}

// Begin implements Handler
func (d *DatabaseHandler) Begin(opts ...*sql.TxOptions) Handler {
	return &DatabaseHandler{db: d.db.Begin(opts...)}
}

// Commit implements Handler
func (d *DatabaseHandler) Commit() error {
	return d.db.Commit().Error
}

// Rollback implements Handler
func (d *DatabaseHandler) Rollback() error {
	return d.db.Rollback().Error
}

// GetDefaultHandler 获取默认的数据库处理器
func GetDefaultHandler() Handler {
	db := Gorm()
	if db == nil {
		return nil
	}
	return NewHandler(db)
}
