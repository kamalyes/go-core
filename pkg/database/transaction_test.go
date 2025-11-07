/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 12:00:00
 * @FilePath: \go-core\pkg\database\transaction_test.go
 * @Description: database 事务操作测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

const (
	UsernameQuery     = "username = ?"
	NestedProductName = "Nested Product"
)

// TransactionTestSuite 事务测试套件
type TransactionTestSuite struct {
	suite.Suite
	db      *gorm.DB
	handler Handler
}

// SetupSuite 测试套件初始化
func (suite *TransactionTestSuite) SetupSuite() {
	db, handler, err := setupTestDB()
	suite.Require().NoError(err)

	suite.db = db
	suite.handler = handler
}

// TearDownSuite 测试套件清理
func (suite *TransactionTestSuite) TearDownSuite() {
	if suite.handler != nil {
		suite.handler.Close()
	}
}

// TestTransactionCommit 测试事务提交
func (suite *TransactionTestSuite) TestTransactionCommit() {
	tx := suite.handler.Begin()
	suite.NotNil(tx)

	user := TestUser{
		Username:   "tx_user_commit",
		Email:      "tx_commit@test.com",
		Age:        25,
		BusinessID: 999,
		Status:     1,
	}

	result := tx.DB().Create(&user)
	suite.NoError(result.Error)
	suite.NotZero(user.ID)

	// 提交事务
	err := tx.Commit()
	suite.NoError(err)

	// 验证数据已提交
	var savedUser TestUser
	suite.handler.DB().Where(UsernameQuery, "tx_user_commit").First(&savedUser)
	suite.Equal("tx_user_commit", savedUser.Username)
}

// TestTransactionRollback 测试事务回滚
func (suite *TransactionTestSuite) TestTransactionRollback() {
	tx := suite.handler.Begin()
	suite.NotNil(tx)

	user := TestUser{
		Username:   "tx_user_rollback",
		Email:      "tx_rollback@test.com",
		Age:        25,
		BusinessID: 999,
		Status:     1,
	}

	result := tx.DB().Create(&user)
	suite.NoError(result.Error)

	// 回滚事务
	err := tx.Rollback()
	suite.NoError(err)

	// 验证数据未保存
	var count int64
	suite.handler.DB().Model(&TestUser{}).Where(UsernameQuery, "tx_user_rollback").Count(&count)
	suite.Equal(int64(0), count)
}

// TestTransactionQuery 测试事务中的查询
func (suite *TransactionTestSuite) TestTransactionQuery() {
	// 先插入一些测试数据
	err := seedTestData(suite.db)
	suite.NoError(err)

	tx := suite.handler.Begin()
	suite.NotNil(tx)

	// 在事务中执行查询
	param := NewSimpleQueryParam(BusinessIDQuery, 1)
	var users []TestUser
	result := tx.Query(param).Find(&users)

	suite.NoError(result.Error)
	suite.Greater(len(users), 0)

	// 提交事务
	err = tx.Commit()
	suite.NoError(err)
}

// TestTransactionNestedOperations 测试事务中的嵌套操作
func (suite *TransactionTestSuite) TestTransactionNestedOperations() {
	tx := suite.handler.Begin()
	suite.NotNil(tx)

	// 创建用户
	user := TestUser{
		Username:   "nested_user",
		Email:      "nested@test.com",
		Age:        30,
		BusinessID: 100,
		Status:     1,
	}

	result := tx.DB().Create(&user)
	suite.NoError(result.Error)

	// 创建相关产品
	product := TestProduct{
		Name:       NestedProductName,
		Price:      150.0,
		BusinessID: user.BusinessID,
		ShopID:     user.ShopID,
		CategoryID: 1,
		Status:     1,
	}

	result = tx.DB().Create(&product)
	suite.NoError(result.Error)

	// 提交事务
	err := tx.Commit()
	suite.NoError(err)

	// 验证两个对象都被保存
	var savedUser TestUser
	suite.handler.DB().Where(UsernameQuery, "nested_user").First(&savedUser)
	suite.Equal("nested_user", savedUser.Username)

	var savedProduct TestProduct
	suite.handler.DB().Where("name = ?", NestedProductName).First(&savedProduct)
	suite.Equal(NestedProductName, savedProduct.Name)
}

// TestTransactionErrorHandling 测试事务错误处理
func (suite *TransactionTestSuite) TestTransactionErrorHandling() {
	tx := suite.handler.Begin()
	suite.NotNil(tx)

	// 尝试插入重复的唯一键数据（如果有唯一约束）
	user1 := TestUser{
		Username:   "duplicate_user",
		Email:      "duplicate@test.com",
		Age:        25,
		BusinessID: 1,
		Status:     1,
	}

	result := tx.DB().Create(&user1)
	suite.NoError(result.Error)

	// 尝试插入相同用户名的用户
	user2 := TestUser{
		Username:   "duplicate_user", // 可能会引起冲突
		Email:      "duplicate2@test.com",
		Age:        30,
		BusinessID: 1,
		Status:     1,
	}

	_ = tx.DB().Create(&user2)
	// 注意：SQLite 内存数据库可能不会严格执行唯一约束
	// 在实际 MySQL/PostgreSQL 中会出错

	// 回滚事务
	err := tx.Rollback()
	suite.NoError(err)

	// 验证没有数据被保存
	var count int64
	suite.handler.DB().Model(&TestUser{}).Where(UsernameQuery, "duplicate_user").Count(&count)
	suite.Equal(int64(0), count)
}

// TestMultipleTransactions 测试多个序列事务
func (suite *TransactionTestSuite) TestMultipleTransactions() {
	// 事务1 - 创建并提交第一个用户
	tx1 := suite.handler.Begin()
	suite.NotNil(tx1)

	user1 := TestUser{
		Username:   "tx1_user",
		Email:      "tx1@test.com",
		Age:        25,
		BusinessID: 1,
		Status:     1,
	}

	result := tx1.DB().Create(&user1)
	suite.NoError(result.Error)

	err := tx1.Commit()
	suite.NoError(err)

	// 事务2 - 创建并提交第二个用户
	tx2 := suite.handler.Begin()
	suite.NotNil(tx2)

	user2 := TestUser{
		Username:   "tx2_user",
		Email:      "tx2@test.com",
		Age:        30,
		BusinessID: 2,
		Status:     1,
	}

	result = tx2.DB().Create(&user2)
	suite.NoError(result.Error)

	err = tx2.Commit()
	suite.NoError(err)

	// 验证两个用户都被保存
	var count int64
	suite.handler.DB().Model(&TestUser{}).Where("username IN ?", []string{"tx1_user", "tx2_user"}).Count(&count)
	suite.Equal(int64(2), count)
}

// TestTransactionWithQueryBuilder 测试事务中使用查询构建器
func (suite *TransactionTestSuite) TestTransactionWithQueryBuilder() {
	// 清理并插入测试数据
	suite.cleanupData()
	err := seedTestData(suite.db)
	suite.NoError(err)

	tx := suite.handler.Begin()
	suite.NotNil(tx)

	// 使用查询构建器
	param := NewQueryBuilder().
		WithBusinessId(1).
		WhereIn("status", []interface{}{1}).
		WithPagination(5, 0).
		Build()

	var users []TestUser
	result := tx.Query(param).Find(&users)
	suite.NoError(result.Error)
	suite.Greater(len(users), 0)

	// 记录原始年龄
	originalAges := make([]int, len(users))
	for i, user := range users {
		originalAges[i] = user.Age
	}

	// 更新这些用户
	for i := range users {
		users[i].Age = users[i].Age + 1
		tx.DB().Save(&users[i])
	}

	// 提交事务
	err = tx.Commit()
	suite.NoError(err)

	// 验证更新
	var updatedUsers []TestUser
	suite.handler.Query(param).Find(&updatedUsers)
	for i, user := range updatedUsers {
		if i < len(originalAges) {
			suite.Equal(originalAges[i]+1, user.Age)
		}
	}
}

// cleanupData 清理测试数据
func (suite *TransactionTestSuite) cleanupData() {
	suite.db.Exec("DELETE FROM test_users")
	suite.db.Exec("DELETE FROM test_products")
}

// TestTransactionTestSuite 运行事务测试套件
func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
