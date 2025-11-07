/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 12:00:00
 * @FilePath: \go-core\pkg\database\utils_test.go
 * @Description: database utils 测试文件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UtilsTestSuite utils测试套件
type UtilsTestSuite struct {
	suite.Suite
}

// SetupSuite 测试套件初始化
func (suite *UtilsTestSuite) SetupSuite() {
	// 测试套件初始化，如果需要的话
}

// TestBuildListQueryOption 测试BuildListQueryOption函数
func (suite *UtilsTestSuite) TestBuildListQueryOption() {
	// 测试传入nil参数
	result := BuildListQueryOption(nil)
	suite.NotNil(result)
	suite.Equal(10, result.Limit)
	suite.Equal(0, result.Offset)
	suite.Equal("DESC", result.Order)
}

// TestBuildListQueryOptionWithValidValues 测试BuildListQueryOption函数 - 有效值
func (suite *UtilsTestSuite) TestBuildListQueryOptionWithValidValues() {
	option := &FindOptionCommon{
		Limit:  20,
		Offset: 10,
		Order:  "ASC",
		By:     "id",
	}
	
	result := BuildListQueryOption(option)
	suite.Equal(20, result.Limit)
	suite.Equal(10, result.Offset)
	suite.Equal("ASC", result.Order)
	suite.Equal("id", result.By)
}

// TestBuildListQueryOptionWithInvalidValues 测试BuildListQueryOption函数 - 无效值
func (suite *UtilsTestSuite) TestBuildListQueryOptionWithInvalidValues() {
	option := &FindOptionCommon{
		Limit:  0,
		Offset: -5,
		Order:  "INVALID",
	}
	
	result := BuildListQueryOption(option)
	suite.Equal(10, result.Limit)    // 应该设置为默认值10
	suite.Equal(0, result.Offset)    // 应该设置为默认值0
	suite.Equal("DESC", result.Order) // 应该设置为默认值DESC
}

// TestBuildListQueryOptionWithEmptyValues 测试BuildListQueryOption函数 - 空值
func (suite *UtilsTestSuite) TestBuildListQueryOptionWithEmptyValues() {
	option := &FindOptionCommon{}
	
	result := BuildListQueryOption(option)
	suite.Equal(10, result.Limit)
	suite.Equal(0, result.Offset)
	suite.Equal("DESC", result.Order)
}

// TestBuildListQueryOptionBoundaryValues 测试BuildListQueryOption函数 - 边界值
func (suite *UtilsTestSuite) TestBuildListQueryOptionBoundaryValues() {
	testCases := []struct {
		name     string
		input    *FindOptionCommon
		expected *FindOptionCommon
	}{
		{
			name: "Limit为1",
			input: &FindOptionCommon{
				Limit: 1,
			},
			expected: &FindOptionCommon{
				Limit:  1,
				Offset: 0,
				Order:  "DESC",
			},
		},
		{
			name: "Offset为0",
			input: &FindOptionCommon{
				Offset: 0,
			},
			expected: &FindOptionCommon{
				Limit:  10,
				Offset: 0,
				Order:  "DESC",
			},
		},
		{
			name: "Order为ASC",
			input: &FindOptionCommon{
				Order: "ASC",
			},
			expected: &FindOptionCommon{
				Limit:  10,
				Offset: 0,
				Order:  "ASC",
			},
		},
	}
	
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := BuildListQueryOption(tc.input)
			suite.Equal(tc.expected.Limit, result.Limit)
			suite.Equal(tc.expected.Offset, result.Offset)
			suite.Equal(tc.expected.Order, result.Order)
		})
	}
}

// TestNewInFilter 测试NewInFilter函数
func (suite *UtilsTestSuite) TestNewInFilter() {
	field := "status"
	values := []interface{}{1, 2, 3}
	
	result := NewInFilter(field, values)
	
	suite.NotNil(result)
	suite.Equal(field, result.DBField)
	suite.Equal(values, result.Values)
	suite.True(result.ExactMatch)
	suite.False(result.AllRegex)
}

// TestNewInFilterWithEmptyValues 测试NewInFilter函数 - 空值
func (suite *UtilsTestSuite) TestNewInFilterWithEmptyValues() {
	field := "category_id"
	values := []interface{}{}
	
	result := NewInFilter(field, values)
	
	suite.NotNil(result)
	suite.Equal(field, result.DBField)
	suite.Equal(values, result.Values)
	suite.True(result.ExactMatch)
}

// TestNewInFilterWithNilValues 测试NewInFilter函数 - nil值
func (suite *UtilsTestSuite) TestNewInFilterWithNilValues() {
	field := "user_id"
	var values []interface{}
	
	result := NewInFilter(field, values)
	
	suite.NotNil(result)
	suite.Equal(field, result.DBField)
	suite.Nil(result.Values)
	suite.True(result.ExactMatch)
}

// TestNewInFilterWithMixedTypes 测试NewInFilter函数 - 混合类型
func (suite *UtilsTestSuite) TestNewInFilterWithMixedTypes() {
	field := "mixed_field"
	values := []interface{}{1, "string", 3.14, true}
	
	result := NewInFilter(field, values)
	
	suite.NotNil(result)
	suite.Equal(field, result.DBField)
	suite.Equal(values, result.Values)
	suite.True(result.ExactMatch)
	suite.Len(result.Values, 4)
}

// TestNewLikeFilter 测试NewLikeFilter函数
func (suite *UtilsTestSuite) TestNewLikeFilter() {
	field := "username"
	values := []interface{}{"admin", "user"}
	allRegex := false
	
	result := NewLikeFilter(field, values, allRegex)
	
	suite.NotNil(result)
	suite.Equal(field, result.DBField)
	suite.Equal(values, result.Values)
	suite.False(result.ExactMatch)
	suite.Equal(allRegex, result.AllRegex)
}

// TestNewLikeFilterWithRegex 测试NewLikeFilter函数 - 使用正则
func (suite *UtilsTestSuite) TestNewLikeFilterWithRegex() {
	field := "email"
	values := []interface{}{".*@gmail.com", ".*@yahoo.com"}
	allRegex := true
	
	result := NewLikeFilter(field, values, allRegex)
	
	suite.NotNil(result)
	suite.Equal(field, result.DBField)
	suite.Equal(values, result.Values)
	suite.False(result.ExactMatch)
	suite.True(result.AllRegex)
}

// TestNewLikeFilterWithEmptyValues 测试NewLikeFilter函数 - 空值
func (suite *UtilsTestSuite) TestNewLikeFilterWithEmptyValues() {
	field := "description"
	values := []interface{}{}
	allRegex := false
	
	result := NewLikeFilter(field, values, allRegex)
	
	suite.NotNil(result)
	suite.Equal(field, result.DBField)
	suite.Equal(values, result.Values)
	suite.False(result.ExactMatch)
	suite.False(result.AllRegex)
}

// TestNewLikeFilterEdgeCases 测试NewLikeFilter函数 - 边界情况
func (suite *UtilsTestSuite) TestNewLikeFilterEdgeCases() {
	testCases := []struct {
		name     string
		field    string
		values   []interface{}
		allRegex bool
	}{
		{
			name:     "空字段名",
			field:    "",
			values:   []interface{}{"test"},
			allRegex: false,
		},
		{
			name:     "单个值",
			field:    "name",
			values:   []interface{}{"single"},
			allRegex: false,
		},
		{
			name:     "多个值",
			field:    "tags",
			values:   []interface{}{"tag1", "tag2", "tag3"},
			allRegex: true,
		},
	}
	
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := NewLikeFilter(tc.field, tc.values, tc.allRegex)
			
			suite.NotNil(result)
			suite.Equal(tc.field, result.DBField)
			suite.Equal(tc.values, result.Values)
			suite.False(result.ExactMatch)
			suite.Equal(tc.allRegex, result.AllRegex)
		})
	}
}

// TestUtilsTestSuite 运行utils测试套件
func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

// 以下是使用testify/assert的单元测试函数，作为suite测试的补充

// TestBuildListQueryOptionAssert 使用assert进行的BuildListQueryOption测试
func TestBuildListQueryOptionAssert(t *testing.T) {
	// 测试nil输入
	result := BuildListQueryOption(nil)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
	assert.Equal(t, "DESC", result.Order)
	
	// 测试有效输入
	option := &FindOptionCommon{
		Limit:  50,
		Offset: 20,
		Order:  "ASC",
		By:     "created_at",
	}
	result = BuildListQueryOption(option)
	assert.Equal(t, 50, result.Limit)
	assert.Equal(t, 20, result.Offset)
	assert.Equal(t, "ASC", result.Order)
	assert.Equal(t, "created_at", result.By)
	
	// 测试无效输入
	invalidOption := &FindOptionCommon{
		Limit:  -1,
		Offset: -10,
		Order:  "UNKNOWN",
	}
	result = BuildListQueryOption(invalidOption)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
	assert.Equal(t, "DESC", result.Order)
}

// TestNewInFilterAssert 使用assert进行的NewInFilter测试
func TestNewInFilterAssert(t *testing.T) {
	field := "id"
	values := []interface{}{1, 2, 3, 4, 5}
	
	result := NewInFilter(field, values)
	
	assert.NotNil(t, result)
	assert.Equal(t, field, result.DBField)
	assert.Equal(t, values, result.Values)
	assert.True(t, result.ExactMatch)
	assert.False(t, result.AllRegex)
	
	// 测试字符串值
	stringValues := []interface{}{"active", "inactive", "pending"}
	result = NewInFilter("status", stringValues)
	assert.Equal(t, stringValues, result.Values)
	assert.True(t, result.ExactMatch)
}

// TestNewLikeFilterAssert 使用assert进行的NewLikeFilter测试
func TestNewLikeFilterAssert(t *testing.T) {
	field := "name"
	values := []interface{}{"john", "jane"}
	
	// 测试非正则模式
	result := NewLikeFilter(field, values, false)
	assert.NotNil(t, result)
	assert.Equal(t, field, result.DBField)
	assert.Equal(t, values, result.Values)
	assert.False(t, result.ExactMatch)
	assert.False(t, result.AllRegex)
	
	// 测试正则模式
	regexValues := []interface{}{"^admin.*", ".*user$"}
	result = NewLikeFilter("username", regexValues, true)
	assert.Equal(t, regexValues, result.Values)
	assert.False(t, result.ExactMatch)
	assert.True(t, result.AllRegex)
}

// BenchmarkBuildListQueryOption 性能基准测试
func BenchmarkBuildListQueryOption(b *testing.B) {
	option := &FindOptionCommon{
		Limit:  20,
		Offset: 10,
		Order:  "ASC",
		By:     "id",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildListQueryOption(option)
	}
}

// BenchmarkBuildListQueryOptionNil 性能基准测试 - nil输入
func BenchmarkBuildListQueryOptionNil(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildListQueryOption(nil)
	}
}

// BenchmarkNewInFilter 性能基准测试
func BenchmarkNewInFilter(b *testing.B) {
	field := "status"
	values := []interface{}{1, 2, 3, 4, 5}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewInFilter(field, values)
	}
}

// BenchmarkNewLikeFilter 性能基准测试
func BenchmarkNewLikeFilter(b *testing.B) {
	field := "name"
	values := []interface{}{"test1", "test2", "test3"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewLikeFilter(field, values, false)
	}
}