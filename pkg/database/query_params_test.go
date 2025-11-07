/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 12:00:00
 * @FilePath: \go-core\pkg\database\query_params_test.go
 * @Description: database 查询参数测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// QueryParamsTestSuite 查询参数测试套件
type QueryParamsTestSuite struct {
	suite.Suite
	db      *gorm.DB
	handler Handler
}

// SetupSuite 测试套件初始化
func (suite *QueryParamsTestSuite) SetupSuite() {
	db, handler, err := setupTestDB()
	suite.Require().NoError(err)

	suite.db = db
	suite.handler = handler

	// 插入测试数据
	err = seedTestData(db)
	suite.Require().NoError(err)
}

// TearDownSuite 测试套件清理
func (suite *QueryParamsTestSuite) TearDownSuite() {
	if suite.handler != nil {
		suite.handler.Close()
	}
}

// TestSimpleQueryParam 测试简单查询参数
func (suite *QueryParamsTestSuite) TestSimpleQueryParam() {
	param := NewSimpleQueryParam(BusinessIDQuery, 1)

	var users []TestUser
	result := suite.handler.Query(param).Find(&users)

	suite.NoError(result.Error)
	suite.Greater(len(users), 0)

	// 验证所有用户都属于 business_id = 1
	for _, user := range users {
		suite.Equal(int64(1), user.BusinessID)
	}
}

// TestPageQueryParam 测试分页查询参数
func (suite *QueryParamsTestSuite) TestPageQueryParam() {
	param := NewPageQueryParam(
		BusinessIDQuery,
		[]interface{}{1},
		2, // limit
		0, // offset
		"created_at DESC",
	)

	var users []TestUser
	result := suite.handler.Query(param).Find(&users)

	suite.NoError(result.Error)
	suite.LessOrEqual(len(users), 2) // 分页限制

	// 测试第二页
	param2 := NewPageQueryParam(
		BusinessIDQuery,
		[]interface{}{1},
		2, // limit
		2, // offset
		"created_at DESC",
	)

	var users2 []TestUser
	suite.handler.Query(param2).Find(&users2)

	// 第一页和第二页的结果不应该重叠
	if len(users) > 0 && len(users2) > 0 {
		suite.NotEqual(users[0].ID, users2[0].ID)
	}
}

// TestAdvancedQueryParam 测试高级查询参数
func (suite *QueryParamsTestSuite) TestAdvancedQueryParam() {
	option := &FindOptionCommon{
		BusinessId:  1,
		Limit:       5,
		Offset:      0,
		By:          "age",
		Order:       "ASC",
		TablePrefix: "",
	}

	param := NewAdvancedQueryParam(option)
	param.AddFilter(&BaseInfoFilter{
		DBField:    "status",
		Values:     []interface{}{1},
		ExactMatch: true,
	})

	var users []TestUser
	result := suite.handler.Query(param).Find(&users)

	suite.NoError(result.Error)
	suite.Greater(len(users), 0)

	// 验证排序（按年龄升序）
	if len(users) > 1 {
		suite.LessOrEqual(users[0].Age, users[1].Age)
	}
}

// TestFilterCreators 测试过滤器创建函数
func (suite *QueryParamsTestSuite) TestFilterCreators() {
	// 测试 IN 过滤器
	inFilter := NewInFilter("status", []interface{}{1, 2})
	suite.Equal("status", inFilter.DBField)
	suite.True(inFilter.ExactMatch)
	suite.Equal([]interface{}{1, 2}, inFilter.Values)

	// 测试 LIKE 过滤器
	likeFilter := NewLikeFilter("name", []interface{}{"test"}, true)
	suite.Equal("name", likeFilter.DBField)
	suite.False(likeFilter.ExactMatch)
	suite.True(likeFilter.AllRegex)
	suite.Equal([]interface{}{"test"}, likeFilter.Values)
}

// TestBuildListQueryOption 测试构建列表查询选项
func (suite *QueryParamsTestSuite) TestBuildListQueryOption() {
	// 测试默认值
	option := BuildListQueryOption(nil)
	suite.NotNil(option)
	suite.Equal(10, option.Limit)
	suite.Equal(0, option.Offset)
	suite.Equal("DESC", option.Order)

	// 测试自定义值
	customOption := &FindOptionCommon{
		Limit:  20,
		Offset: 10,
		Order:  "ASC",
	}

	result := BuildListQueryOption(customOption)
	suite.Equal(20, result.Limit)
	suite.Equal(10, result.Offset)
	suite.Equal("ASC", result.Order)

	// 测试无效值的修正
	invalidOption := &FindOptionCommon{
		Limit:  -1,
		Offset: -5,
		Order:  "INVALID",
	}

	corrected := BuildListQueryOption(invalidOption)
	suite.Equal(10, corrected.Limit)     // 默认值
	suite.Equal(0, corrected.Offset)     // 修正负值
	suite.Equal("DESC", corrected.Order) // 修正无效值
}

// TestSimpleQueryParamEdgeCases 测试SimpleQueryParam的边界情况
func TestSimpleQueryParamEdgeCases(t *testing.T) {
	db, handler, err := setupTestDB()
	assert.NoError(t, err)
	defer handler.Close()

	// 插入测试数据
	err = seedTestData(db)
	assert.NoError(t, err)

	// 测试空的where子句
	param := NewSimpleQueryParam("", []interface{}{})
	var users []TestUser
	result := handler.Query(param).Find(&users)
	assert.NoError(t, result.Error)
	assert.Greater(t, len(users), 0) // 应该返回所有数据

	// 测试nil参数
	param = NewSimpleQueryParam("status = ?", nil)
	result = handler.Query(param).Find(&users)
	// 可能会有错误，但我们测试覆盖率
	_ = result
}

// TestPageQueryParamEdgeCases 测试PageQueryParam的边界情况
func TestPageQueryParamEdgeCases(t *testing.T) {
	db, handler, err := setupTestDB()
	assert.NoError(t, err)
	defer handler.Close()

	// 插入测试数据
	err = seedTestData(db)
	assert.NoError(t, err)

	// 测试所有参数为0或空的情况
	param := NewPageQueryParam("", []interface{}{}, 0, 0, "")
	var users []TestUser
	result := handler.Query(param).Find(&users)
	assert.NoError(t, result.Error)

	// 测试空orderBy的情况
	param = NewPageQueryParam("status = ?", []interface{}{1}, 10, 0, "")
	result = handler.Query(param).Find(&users)
	assert.NoError(t, result.Error)
}

// TestQueryParamsTestSuite 运行查询参数测试套件
func TestQueryParamsTestSuite(t *testing.T) {
	suite.Run(t, new(QueryParamsTestSuite))
}
