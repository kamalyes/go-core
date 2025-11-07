/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 10:11:55
 * @FilePath: \go-core\pkg\database\integration_test.go
 * @Description: database 集成测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// IntegrationTestSuite 集成测试套件
type IntegrationTestSuite struct {
	suite.Suite
	db      *gorm.DB
	handler Handler
}

// SetupSuite 测试套件初始化
func (suite *IntegrationTestSuite) SetupSuite() {
	db, handler, err := setupTestDB()
	suite.Require().NoError(err)

	suite.db = db
	suite.handler = handler

	// 插入测试数据
	err = seedTestData(db)
	suite.Require().NoError(err)
}

// TearDownSuite 测试套件清理
func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.handler != nil {
		suite.handler.Close()
	}
}

// TestCompleteUserWorkflow 测试完整的用户操作流程
func (suite *IntegrationTestSuite) TestCompleteUserWorkflow() {
	// 1. 创建新用户
	tx := suite.handler.Begin()
	user := TestUser{
		Username:   "workflow_user",
		Email:      "workflow@test.com",
		Age:        25,
		BusinessID: 1,
		ShopID:     101,
		Status:     1,
		Tags:       "new,active",
	}

	result := tx.DB().Create(&user)
	suite.NoError(result.Error)
	suite.NotZero(user.ID)

	err := tx.Commit()
	suite.NoError(err)

	// 2. 查询用户
	param := NewQueryBuilder().
		WithBusinessId(1).
		WhereIn("status", []interface{}{1}).
		WhereFindInSet("tags", []string{"active"}).
		Build()

	var users []TestUser
	result = suite.handler.Query(param).Find(&users)
	suite.NoError(result.Error)
	suite.Greater(len(users), 0)

	// 验证新创建的用户在结果中
	found := false
	for _, u := range users {
		if u.Username == "workflow_user" {
			found = true
			break
		}
	}
	suite.True(found)

	// 3. 更新用户信息
	tx = suite.handler.Begin()
	result = tx.DB().Model(&user).Update("age", 26)
	suite.NoError(result.Error)
	err = tx.Commit()
	suite.NoError(err)

	// 4. 验证更新
	var updatedUser TestUser
	suite.handler.DB().Where(UsernameQuery, "workflow_user").First(&updatedUser)
	suite.Equal(26, updatedUser.Age)

	// 5. 删除用户
	tx = suite.handler.Begin()
	result = tx.DB().Delete(&user)
	suite.NoError(result.Error)
	err = tx.Commit()
	suite.NoError(err)

	// 6. 验证删除
	var count int64
	suite.handler.DB().Model(&TestUser{}).Where(UsernameQuery, "workflow_user").Count(&count)
	suite.Equal(int64(0), count)
}

// TestBusinessScenario 测试业务场景
func (suite *IntegrationTestSuite) TestBusinessScenario() {
	// 场景：查询特定业务下的活跃用户，并创建相关产品

	// 1. 查询业务1下的活跃用户
	param := NewQueryBuilder().
		WithBusinessId(1).
		WhereIn("status", []interface{}{1}).
		WithOrder("age", "ASC").
		WithPagination(10, 0).
		Build()

	var activeUsers []TestUser
	result := suite.handler.Query(param).Find(&activeUsers)
	suite.NoError(result.Error)
	suite.Greater(len(activeUsers), 0)

	// 2. 为每个用户创建相关产品
	tx := suite.handler.Begin()

	for _, user := range activeUsers {
		product := TestProduct{
			Name:       "Product for " + user.Username,
			Price:      100.0 + float64(user.Age),
			BusinessID: user.BusinessID,
			ShopID:     user.ShopID,
			CategoryID: 1,
			Status:     1,
		}

		result = tx.DB().Create(&product)
		suite.NoError(result.Error)
	}

	err := tx.Commit()
	suite.NoError(err)

	// 3. 验证产品创建
	var productCount int64
	suite.handler.DB().Model(&TestProduct{}).Where(BusinessIDQuery, 1).Count(&productCount)
	suite.GreaterOrEqual(productCount, int64(len(activeUsers)))

	// 4. 查询产品统计
	advancedParam := NewAdvancedQueryParam(&FindOptionCommon{
		BusinessId: 1,
		GroupBy:    "shop_id",
		By:         "shop_id",
		Order:      "ASC",
	})

	var shopProducts []TestProduct
	suite.handler.Query(advancedParam).Find(&shopProducts)
	suite.Greater(len(shopProducts), 0)
}

// TestConcurrentOperations 测试序列化操作（SQLite并发写入限制）
func (suite *IntegrationTestSuite) TestConcurrentOperations() {
	// 由于SQLite的并发写入限制，改为序列化测试
	const operations = 5
	
	for i := 0; i < operations; i++ {
		tx := suite.handler.Begin()

		user := TestUser{
			Username:   "concurrent_user_" + string(rune('0'+i)),
			Email:      "concurrent" + string(rune('0'+i)) + "@test.com",
			Age:        20 + i,
			BusinessID: int64(1 + (i % 2)),
			ShopID:     int64(100 + i),
			Status:     1,
		}

		result := tx.DB().Create(&user)
		suite.NoError(result.Error)

		err := tx.Commit()
		suite.NoError(err)
	}

	// 验证所有用户都被创建
	var count int64
	suite.handler.DB().Model(&TestUser{}).Where("username LIKE ?", "concurrent_user_%").Count(&count)
	suite.Equal(int64(operations), count)
}

// TestComplexQueryScenario 测试复杂查询场景
func (suite *IntegrationTestSuite) TestComplexQueryScenario() {
	// 场景：多条件查询 + 时间范围 + 分页 + 排序

	param := NewQueryBuilder().
		WithBusinessId(1).
		WhereIn("status", []interface{}{1, 2}).
		WhereIn("age", []interface{}{25, 30, 35}).
		WhereFindInSet("tags", []string{"active", "vip"}).
		WithOrder("created_at", "DESC").
		WithPagination(5, 0).
		Build()

	var users []TestUser
	result := suite.handler.Query(param).Find(&users)
	suite.NoError(result.Error)

	// 验证查询结果符合条件
	for _, user := range users {
		suite.Equal(int64(1), user.BusinessID)
		suite.Contains([]int{1, 2}, user.Status)
		suite.Contains([]int{25, 30, 35}, user.Age)
	}

	// 测试总数查询
	var total int64
	countParam := NewQueryBuilder().
		WithBusinessId(1).
		WhereIn("status", []interface{}{1, 2}).
		WhereIn("age", []interface{}{25, 30, 35}).
		Build()

	suite.handler.Query(countParam).Model(&TestUser{}).Count(&total)
	suite.GreaterOrEqual(total, int64(len(users)))
}

// TestErrorHandlingScenario 测试错误处理场景
func (suite *IntegrationTestSuite) TestErrorHandlingScenario() {
	// 测试事务回滚场景
	tx := suite.handler.Begin()

	// 创建一个用户
	user := TestUser{
		Username:   "error_user",
		Email:      "error@test.com",
		Age:        25,
		BusinessID: 999,
		Status:     1,
	}

	result := tx.DB().Create(&user)
	suite.NoError(result.Error)

	// 模拟业务错误，需要回滚
	simulateError := true
	if simulateError {
		err := tx.Rollback()
		suite.NoError(err)

		// 验证用户没有被创建
		var count int64
		suite.handler.DB().Model(&TestUser{}).Where(UsernameQuery, "error_user").Count(&count)
		suite.Equal(int64(0), count)
	}
}

// TestDataConsistency 测试数据一致性
func (suite *IntegrationTestSuite) TestDataConsistency() {
	// 创建用户和相关产品，确保数据一致性
	tx := suite.handler.Begin()

	user := TestUser{
		Username:   "consistency_user",
		Email:      "consistency@test.com",
		Age:        30,
		BusinessID: 100,
		ShopID:     200,
		Status:     1,
	}

	result := tx.DB().Create(&user)
	suite.NoError(result.Error)

	// 创建相关产品
	products := []TestProduct{
		{
			Name:       "Product 1 for consistency_user",
			Price:      100.0,
			BusinessID: user.BusinessID,
			ShopID:     user.ShopID,
			CategoryID: 1,
			Status:     1,
		},
		{
			Name:       "Product 2 for consistency_user",
			Price:      200.0,
			BusinessID: user.BusinessID,
			ShopID:     user.ShopID,
			CategoryID: 2,
			Status:     1,
		},
	}

	for _, product := range products {
		result = tx.DB().Create(&product)
		suite.NoError(result.Error)
	}

	err := tx.Commit()
	suite.NoError(err)

	// 验证数据一致性
	var userCount int64
	suite.handler.DB().Model(&TestUser{}).Where("business_id = ? AND shop_id = ?", 100, 200).Count(&userCount)
	suite.Equal(int64(1), userCount)

	var productCount int64
	suite.handler.DB().Model(&TestProduct{}).Where("business_id = ? AND shop_id = ?", 100, 200).Count(&productCount)
	suite.Equal(int64(2), productCount)
}

// TestIntegrationTestSuite 运行集成测试套件
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
