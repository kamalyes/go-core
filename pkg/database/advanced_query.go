/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 10:33:02
 * @FilePath: \go-core\pkg\database\advanced_query.go
 * @Description: 高级查询参数实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"strings"

	"gorm.io/gorm"
)

// 数据库方言常量
const (
	// SQL模式常量
	SQLLikePattern   = " LIKE ?"
	SQLEqualsPattern = " = ?"
)

// AdvancedQueryParam 高级查询参数，支持复杂查询构建
type AdvancedQueryParam struct {
	option     *FindOptionCommon
	filters    []*BaseInfoFilter
	timeRanges map[string][2]string // 时间范围查询 key: 字段名, value: [开始时间, 结束时间]
	findInSets map[string][]string  // FIND_IN_SET查询 key: 字段名, value: 查找值列表
}

// NewAdvancedQueryParam 创建高级查询参数
func NewAdvancedQueryParam(option *FindOptionCommon) *AdvancedQueryParam {
	if option == nil {
		option = &FindOptionCommon{}
	}
	return &AdvancedQueryParam{
		option:     option,
		filters:    make([]*BaseInfoFilter, 0),
		timeRanges: make(map[string][2]string),
		findInSets: make(map[string][]string),
	}
}

// AddFilter 添加过滤器
func (a *AdvancedQueryParam) AddFilter(filter *BaseInfoFilter) *AdvancedQueryParam {
	if filter != nil {
		a.filters = append(a.filters, filter)
	}
	return a
}

// AddTimeRange 添加时间范围查询
func (a *AdvancedQueryParam) AddTimeRange(field, startTime, endTime string) *AdvancedQueryParam {
	if field != "" && startTime != "" && endTime != "" {
		a.timeRanges[field] = [2]string{startTime, endTime}
	}
	return a
}

// AddFindInSet 添加FIND_IN_SET查询
func (a *AdvancedQueryParam) AddFindInSet(field string, values []string) *AdvancedQueryParam {
	if field != "" && len(values) > 0 {
		a.findInSets[field] = values
	}
	return a
}

// Where 实现 QueryParam 接口
func (a *AdvancedQueryParam) Where(db *gorm.DB) *gorm.DB {
	db = a.applyBusinessAndShopConditions(db)
	db = a.applyFilters(db)
	db = a.applyTimeRangeConditions(db)
	db = a.applyFindInSetConditions(db)
	db = a.applyGroupAndOrder(db)
	db = a.applyPagination(db)
	return db
}

// applyBusinessAndShopConditions 应用业务和店铺条件
func (a *AdvancedQueryParam) applyBusinessAndShopConditions(db *gorm.DB) *gorm.DB {
	if a.option.ExcludeBusinessAndShop {
		return db
	}

	if !a.option.ExcludeBusiness && (a.option.BusinessId > 0 || a.option.IncludeBusinessIdZero) {
		fieldName := a.option.TablePrefix + "business_id"
		db = db.Where(fieldName+" = ?", a.option.BusinessId)
	}

	if !a.option.ExcludeShop && a.option.ShopId > 0 {
		fieldName := a.option.TablePrefix + "shop_id"
		db = db.Where(fieldName+" = ?", a.option.ShopId)
	}

	return db
}

// applyFilters 应用过滤器
func (a *AdvancedQueryParam) applyFilters(db *gorm.DB) *gorm.DB {
	for _, filter := range a.filters {
		db = a.applyFilter(db, filter)
	}
	return db
}

// applyFilter 应用单个过滤器
func (a *AdvancedQueryParam) applyFilter(db *gorm.DB, filter *BaseInfoFilter) *gorm.DB {
	if len(filter.Values) == 0 || filter.DBField == "" {
		return db
	}

	if filter.ExactMatch {
		return db.Where(filter.DBField+" IN (?)", filter.Values)
	}

	// 模糊匹配处理 - 使用LIKE而不是REGEXP以兼容SQLite
	var conditions []string
	var args []interface{}

	for _, valueItem := range filter.Values {
		if str, ok := valueItem.(string); ok {
			if filter.AllRegex {
				// 全模匹配：%value%
				conditions = append(conditions, filter.DBField+" LIKE ?")
				args = append(args, "%"+str+"%")
			} else {
				// 左模匹配：value%
				conditions = append(conditions, filter.DBField+" LIKE ?")
				args = append(args, str+"%")
			}
		}
	}

	if len(conditions) > 0 {
		whereClause := "(" + strings.Join(conditions, " OR ") + ")"
		return db.Where(whereClause, args...)
	}

	return db
}

// applyTimeRangeConditions 应用时间范围条件
func (a *AdvancedQueryParam) applyTimeRangeConditions(db *gorm.DB) *gorm.DB {
	for field, timeRange := range a.timeRanges {
		db = db.Where(field+" BETWEEN ? AND ?", timeRange[0], timeRange[1])
	}
	return db
}

// applyFindInSetConditions 应用FIND_IN_SET条件
func (a *AdvancedQueryParam) applyFindInSetConditions(db *gorm.DB) *gorm.DB {
	for field, values := range a.findInSets {
		if len(values) == 1 {
			db = a.buildFindInSetCondition(db, field, values[0])
		} else if len(values) > 1 {
			var conditions []string
			var args []interface{}
			for _, value := range values {
				// 检测数据库类型并使用相应的语法
				dialectName := db.Dialector.Name()
				if dialectName == "sqlite" {
					// SQLite 使用 LIKE 和 % 通配符来模拟 FIND_IN_SET
					condition := a.buildSQLiteFindInSetCondition(field)
					conditions = append(conditions, condition)
					args = append(args, value+",%", "%,"+value+",%", "%,"+value, value)
				} else {
					// MySQL 等其他数据库使用 FIND_IN_SET
					conditions = append(conditions, "FIND_IN_SET(?, "+field+")")
					args = append(args, value)
				}
			}
			whereClause := "(" + strings.Join(conditions, " OR ") + ")"
			db = db.Where(whereClause, args...)
		}
	}
	return db
}

// buildSQLiteFindInSetCondition 构建SQLite的FIND_IN_SET条件
func (a *AdvancedQueryParam) buildSQLiteFindInSetCondition(field string) string {
	// 构建SQLite的FIND_IN_SET替代语法
	// 处理四种情况: "value", "value,xxx", "xxx,value", "xxx,value,yyy"
	patterns := []string{
		field + SQLLikePattern,   // value,%
		field + SQLLikePattern,   // %,value,%
		field + SQLLikePattern,   // %,value
		field + SQLEqualsPattern, // value
	}
	return "(" + strings.Join(patterns, " OR ") + ")"
}

// buildFindInSetCondition 构建单个FIND_IN_SET条件
func (a *AdvancedQueryParam) buildFindInSetCondition(db *gorm.DB, field string, value string) *gorm.DB {
	dialectName := db.Dialector.Name()
	if dialectName == "sqlite" {
		// SQLite 使用 LIKE 模式匹配
		condition := a.buildSQLiteFindInSetCondition(field)
		return db.Where(condition, value+",%", "%,"+value+",%", "%,"+value, value)
	}
	// MySQL 等其他数据库使用 FIND_IN_SET
	return db.Where("FIND_IN_SET(?, "+field+")", value)
}

// applyGroupAndOrder 应用分组和排序
func (a *AdvancedQueryParam) applyGroupAndOrder(db *gorm.DB) *gorm.DB {
	if a.option.GroupBy != "" {
		db = db.Group(a.option.GroupBy)
	}

	if !a.option.DisableOrderBy && a.option.By != "" {
		orderField := a.option.TablePrefix + a.option.By
		orderDirection := "DESC"
		if a.option.Order != "" {
			orderDirection = a.option.Order
		}
		db = db.Order(orderField + " " + orderDirection)
	}

	return db
}

// applyPagination 应用分页
func (a *AdvancedQueryParam) applyPagination(db *gorm.DB) *gorm.DB {
	if a.option.Limit > 0 {
		db = db.Limit(a.option.Limit)
	}
	if a.option.Offset > 0 {
		db = db.Offset(a.option.Offset)
	}
	return db
}
