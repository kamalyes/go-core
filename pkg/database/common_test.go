/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 10:19:26
 * @FilePath: \go-core\pkg\database\common_test.go
 * @Description: database 测试公共定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	InMemoryDB      = "file::memory:?cache=shared"
	BusinessIDQuery = "business_id = ?"
)

// TestUser 测试用户模型
type TestUser struct {
	ID         uint   `gorm:"primarykey"`
	Username   string `gorm:"size:50;not null"`
	Email      string `gorm:"size:100;not null"`
	Age        int
	BusinessID int64
	ShopID     int64
	Status     int
	Tags       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName 指定表名
func (TestUser) TableName() string {
	return "test_users"
}

// TestProduct 测试产品模型
type TestProduct struct {
	ID         uint   `gorm:"primarykey"`
	Name       string `gorm:"size:100;not null"`
	Price      float64
	BusinessID int64
	ShopID     int64
	CategoryID int
	Status     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName 指定表名
func (TestProduct) TableName() string {
	return "test_products"
}

// setupTestDB 创建测试数据库
func setupTestDB() (*gorm.DB, Handler, error) {
	db, err := gorm.Open(sqlite.Open(InMemoryDB), &gorm.Config{
		// 启用日志以便调试
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}

	// 确保连接是活跃的
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	
	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return nil, nil, err
	}

	handler := NewHandler(db)

	// 自动迁移测试表
	err = handler.AutoMigrate(&TestUser{}, &TestProduct{})
	if err != nil {
		return nil, nil, err
	}

	// 验证表已创建
	if !db.Migrator().HasTable(&TestUser{}) {
		return nil, nil, fmt.Errorf("test_users table not created")
	}
	
	if !db.Migrator().HasTable(&TestProduct{}) {
		return nil, nil, fmt.Errorf("test_products table not created")
	}

	return db, handler, nil
}

// seedTestData 插入测试数据
func seedTestData(db *gorm.DB) error {
	users := []TestUser{
		{Username: "john_doe", Email: "john@test.com", Age: 25, BusinessID: 1, ShopID: 101, Status: 1, Tags: "vip,active"},
		{Username: "jane_smith", Email: "jane@test.com", Age: 30, BusinessID: 1, ShopID: 102, Status: 1, Tags: "active,premium"},
		{Username: "bob_wilson", Email: "bob@test.com", Age: 35, BusinessID: 2, ShopID: 201, Status: 2, Tags: "inactive"},
		{Username: "alice_brown", Email: "alice@test.com", Age: 28, BusinessID: 1, ShopID: 101, Status: 1, Tags: "vip"},
		{Username: "charlie_davis", Email: "charlie@test.com", Age: 32, BusinessID: 2, ShopID: 202, Status: 1, Tags: "active"},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	products := []TestProduct{
		{Name: "Product A", Price: 100.0, BusinessID: 1, ShopID: 101, CategoryID: 1, Status: 1},
		{Name: "Product B", Price: 200.0, BusinessID: 1, ShopID: 102, CategoryID: 2, Status: 1},
		{Name: "Product C", Price: 300.0, BusinessID: 2, ShopID: 201, CategoryID: 1, Status: 2},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
