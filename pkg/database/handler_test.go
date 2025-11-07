/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 12:00:00
 * @FilePath: \go-core\pkg\database\handler_test.go
 * @Description: database Handler 接口测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// HandlerTestSuite Handler 接口测试套件
type HandlerTestSuite struct {
	suite.Suite
	db      *gorm.DB
	handler Handler
}

// SetupSuite 测试套件初始化
func (suite *HandlerTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(InMemoryDB), &gorm.Config{})
	suite.Require().NoError(err)

	suite.db = db
	suite.handler = NewHandler(db)
}

// TearDownSuite 测试套件清理
func (suite *HandlerTestSuite) TearDownSuite() {
	if suite.handler != nil {
		suite.handler.Close()
	}
}

// TestHandlerDB 测试 DB 方法
func (suite *HandlerTestSuite) TestHandlerDB() {
	db := suite.handler.DB()
	suite.NotNil(db)
	suite.Equal(suite.db, db)
}

// TestHandlerAutoMigrate 测试 AutoMigrate 方法
func (suite *HandlerTestSuite) TestHandlerAutoMigrate() {
	type TempModel struct {
		ID   uint   `gorm:"primarykey"`
		Name string `gorm:"size:100"`
	}

	err := suite.handler.AutoMigrate(&TempModel{})
	suite.NoError(err)

	// 验证表是否创建成功
	suite.True(suite.db.Migrator().HasTable(&TempModel{}))
}

// TestHandlerClose 测试关闭数据库连接
func (suite *HandlerTestSuite) TestHandlerClose() {
	// 创建新的处理器进行测试
	db, err := gorm.Open(sqlite.Open(InMemoryDB), &gorm.Config{})
	suite.NoError(err)

	handler := NewHandler(db)

	// 测试关闭连接
	err = handler.Close()
	suite.NoError(err)
}

// TestNewHandler 测试 NewHandler 函数
func TestNewHandler(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(InMemoryDB), &gorm.Config{})
	assert.NoError(t, err)

	handler := NewHandler(db)
	assert.NotNil(t, handler)
	assert.Equal(t, db, handler.DB())
}

// TestGetDefaultHandler 测试 GetDefaultHandler 函数
func TestGetDefaultHandler(t *testing.T) {
	// 这个测试需要实际的数据库配置，在测试环境中跳过
	// 因为它依赖于全局配置变量 global.CONFIG
	t.Skip("Skipping TestGetDefaultHandler: requires global configuration setup")
}

// TestHandlerTestSuite 运行 Handler 测试套件
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
