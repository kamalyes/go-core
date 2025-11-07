/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 09:15:07
 * @FilePath: \go-core\pkg\database\models.go
 * @Description: 数据库查询相关数据模型
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

// BaseInfoFilter 基础过滤器
type BaseInfoFilter struct {
	DBField    string        // DB中字段名称
	Values     []interface{} // values
	ExactMatch bool          // 是否精确查询
	AllRegex   bool          // 标记是否进行全模匹配(默认为false，大部分需要左模匹配命中索引)只在模糊查询时有效
}

// FindOptionCommon 通用查询选项
type FindOptionCommon struct {
	BusinessId               int64  // 业务ID
	ShopId                  int64  // 店铺ID
	ExcludeBusiness         bool   // 排除业务ID条件
	ExcludeShop            bool   // 排除店铺ID条件
	ExcludeBusinessAndShop bool   // 排除业务和店铺条件
	IncludeBusinessIdZero  bool   // 包含业务ID为0的情况
	GroupBy                string // 分组字段
	By                     string // 排序字段
	Order                  string // 排序方向 ASC/DESC
	DisableOrderBy         bool   // 禁用排序
	Limit                  int    // 限制数量
	Offset                 int    // 偏移量
	TablePrefix            string // 表前缀
}