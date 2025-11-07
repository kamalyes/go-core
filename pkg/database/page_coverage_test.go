/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 12:00:00
 * @FilePath: \go-core\pkg\database\page_coverage_test.go
 * @Description: page.go的额外覆盖率测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/kamalyes/go-core/pkg/global"
	"github.com/stretchr/testify/assert"
)

// TestFindPageErrorHandling 测试FindPage的错误处理
func TestFindPageErrorHandling(t *testing.T) {
	var users []TestUser
	var user TestUser

	// 测试pageInfo为nil的情况
	pageResult, err := FindPage(&user, &users, nil)
	assert.Error(t, err)
	assert.Nil(t, pageResult)
	assert.Contains(t, err.Error(), "入参pageInfo不能为空指针")
}

// TestPageInfoValidation 测试PageInfo的验证
func TestPageInfoValidation(t *testing.T) {
	db, handler, err := setupTestDB()
	assert.NoError(t, err)
	defer handler.Close()

	// 设置global.DB以供FindPage使用
	originalDB := global.DB
	global.DB = db
	defer func() {
		global.DB = originalDB
	}()

	// 插入测试数据
	err = seedTestData(db)
	assert.NoError(t, err)

	var users []TestUser
	var user TestUser

	// 测试负数页码
	pageInfo := &PageInfo{
		Current:  -1,
		RowCount: 10,
	}

	pageResult, err := FindPage(&user, &users, pageInfo)
	assert.NoError(t, err)
	assert.NotNil(t, pageResult)
	// 负数页码应该被处理
	assert.Equal(t, -1, pageResult.Page) // 实际上PageInfo.Current直接赋值给PageBean.Page

	// 测试零页码
	pageInfo = &PageInfo{
		Current:  0,
		RowCount: 10,
	}

	pageResult, err = FindPage(&user, &users, pageInfo)
	assert.NoError(t, err)
	assert.NotNil(t, pageResult)
	assert.Equal(t, 0, pageResult.Page)

	// 测试负数页面大小
	pageInfo = &PageInfo{
		Current:  1,
		RowCount: -5,
	}

	pageResult, err = FindPage(&user, &users, pageInfo)
	assert.NoError(t, err)
	assert.NotNil(t, pageResult)
	assert.Equal(t, -5, pageResult.PageSize) // 实际上PageInfo.RowCount直接赋值给PageBean.PageSize

	// 测试零页面大小
	pageInfo = &PageInfo{
		Current:  1,
		RowCount: 0,
	}

	pageResult, err = FindPage(&user, &users, pageInfo)
	assert.NoError(t, err)
	assert.NotNil(t, pageResult)
	assert.Equal(t, 0, pageResult.PageSize)
}

// 跳过MockDB测试，因为它需要更复杂的GORM模拟设置

// TestBufferAppendWithNilValues 测试Buffer处理nil值
func TestBufferAppendWithNilValues(t *testing.T) {
	buffer := NewBuffer()

	// 测试append nil值（不支持的类型不会添加任何内容）
	buffer.Append(nil)
	result := buffer.String()
	assert.Equal(t, "", result) // Buffer的Append方法不处理nil

	// 测试append指向nil的指针
	var nilPtr *string
	buffer = NewBuffer()
	buffer.Append(nilPtr)
	result = buffer.String()
	assert.Equal(t, "", result) // Buffer的Append方法不处理指针类型
}

// TestBufferAppendComplexTypes 测试Buffer处理复杂类型
func TestBufferAppendComplexTypes(t *testing.T) {
	buffer := NewBuffer()

	// 测试切片（不支持的类型）
	slice := []int{1, 2, 3}
	buffer.Append(slice)
	result := buffer.String()
	assert.Equal(t, "", result) // Buffer的Append方法不处理切片

	// 测试map（不支持的类型）
	buffer = NewBuffer()
	testMap := map[string]int{"a": 1, "b": 2}
	buffer.Append(testMap)
	result = buffer.String()
	assert.Equal(t, "", result) // Buffer的Append方法不处理map

	// 测试结构体（不支持的类型）
	buffer = NewBuffer()
	type testStruct struct {
		Name string
		Age  int
	}
	ts := testStruct{Name: "test", Age: 25}
	buffer.Append(ts)
	result = buffer.String()
	assert.Equal(t, "", result) // Buffer的Append方法不处理结构体
}

// TestCamelToCaseEdgeCases 测试CamelToCase的边界情况
func TestCamelToCaseEdgeCases(t *testing.T) {
	// 测试空字符串
	assert.Equal(t, "", CamelToCase(""))

	// 测试单个字符
	assert.Equal(t, "a", CamelToCase("a"))
	assert.Equal(t, "a", CamelToCase("A"))

	// 测试连续大写字母
	assert.Equal(t, "u_r_l", CamelToCase("URL"))
	assert.Equal(t, "h_t_t_p_u_r_l", CamelToCase("HTTPURL"))

	// 测试数字
	assert.Equal(t, "user_i_d2", CamelToCase("userID2"))
	assert.Equal(t, "user2_name", CamelToCase("user2Name"))

	// 测试特殊字符
	assert.Equal(t, "user__name", CamelToCase("user_Name"))
	assert.Equal(t, "user-_name", CamelToCase("user-Name"))

	// 测试已经是snake_case的字符串
	assert.Equal(t, "user_name", CamelToCase("user_name"))

	// 测试混合情况
	assert.Equal(t, "get_u_s_e_r_by_i_d", CamelToCase("getUSERByID"))
}
