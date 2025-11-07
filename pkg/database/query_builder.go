/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 09:15:07
 * @FilePath: \go-core\pkg\database\query_builder.go
 * @Description: 查询构建器实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

// QueryBuilder 查询构建器
type QueryBuilder struct {
	param *AdvancedQueryParam
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		param: NewAdvancedQueryParam(nil),
	}
}

// WithOption 设置查询选项
func (qb *QueryBuilder) WithOption(option *FindOptionCommon) *QueryBuilder {
	qb.param.option = BuildListQueryOption(option)
	return qb
}

// WithBusinessId 设置业务ID
func (qb *QueryBuilder) WithBusinessId(businessId int64) *QueryBuilder {
	qb.param.option.BusinessId = businessId
	return qb
}

// WithShopId 设置店铺ID
func (qb *QueryBuilder) WithShopId(shopId int64) *QueryBuilder {
	qb.param.option.ShopId = shopId
	return qb
}

// WithTablePrefix 设置表前缀
func (qb *QueryBuilder) WithTablePrefix(prefix string) *QueryBuilder {
	qb.param.option.TablePrefix = prefix
	return qb
}

// WithPagination 设置分页
func (qb *QueryBuilder) WithPagination(limit, offset int) *QueryBuilder {
	qb.param.option.Limit = limit
	qb.param.option.Offset = offset
	return qb
}

// WithOrder 设置排序
func (qb *QueryBuilder) WithOrder(field, order string) *QueryBuilder {
	qb.param.option.By = field
	qb.param.option.Order = order
	return qb
}

// WithGroupBy 设置分组
func (qb *QueryBuilder) WithGroupBy(groupBy string) *QueryBuilder {
	qb.param.option.GroupBy = groupBy
	return qb
}

// WhereIn 添加IN条件
func (qb *QueryBuilder) WhereIn(field string, values []interface{}) *QueryBuilder {
	qb.param.AddFilter(NewInFilter(field, values))
	return qb
}

// WhereLike 添加LIKE条件
func (qb *QueryBuilder) WhereLike(field string, values []interface{}, allRegex bool) *QueryBuilder {
	qb.param.AddFilter(NewLikeFilter(field, values, allRegex))
	return qb
}

// WhereTimeRange 添加时间范围条件
func (qb *QueryBuilder) WhereTimeRange(field, startTime, endTime string) *QueryBuilder {
	qb.param.AddTimeRange(field, startTime, endTime)
	return qb
}

// WhereFindInSet 添加FIND_IN_SET条件
func (qb *QueryBuilder) WhereFindInSet(field string, values []string) *QueryBuilder {
	qb.param.AddFindInSet(field, values)
	return qb
}

// Build 构建查询参数
func (qb *QueryBuilder) Build() QueryParam {
	return qb.param
}