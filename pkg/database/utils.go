/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 09:15:07
 * @FilePath: \go-core\pkg\database\utils.go
 * @Description: 数据库查询工具函数
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

// BuildListQueryOption 构建列表查询选项
func BuildListQueryOption(option *FindOptionCommon) *FindOptionCommon {
	if option == nil {
		option = &FindOptionCommon{}
	}

	// 设置默认值
	if option.Limit < 1 {
		option.Limit = 10
	}
	if option.Offset < 0 {
		option.Offset = 0
	}
	if option.Order != "ASC" && option.Order != "DESC" {
		option.Order = "DESC"
	}

	return option
}

// NewInFilter 创建IN查询过滤器
func NewInFilter(field string, values []interface{}) *BaseInfoFilter {
	return &BaseInfoFilter{
		DBField:    field,
		Values:     values,
		ExactMatch: true,
	}
}

// NewLikeFilter 创建LIKE查询过滤器
func NewLikeFilter(field string, values []interface{}, allRegex bool) *BaseInfoFilter {
	return &BaseInfoFilter{
		DBField:    field,
		Values:     values,
		ExactMatch: false,
		AllRegex:   allRegex,
	}
}