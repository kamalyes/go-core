/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 12:00:00
 * @FilePath: \go-core\pkg\database\benchmark_test.go
 * @Description: database 性能基准测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"
	"time"
)

// BenchmarkUser 基准测试用户模型
type BenchmarkUser struct {
	ID         uint   `gorm:"primarykey"`
	Username   string `gorm:"size:50;not null;index"`
	Email      string `gorm:"size:100;not null;index"`
	Age        int    `gorm:"index"`
	BusinessID int64  `gorm:"index"`
	ShopID     int64  `gorm:"index"`
	Status     int    `gorm:"index"`
	Tags       string
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time
}

// setupBenchmarkDB 设置基准测试数据库
func setupBenchmarkDB(b *testing.B) (Handler, func()) {
	_, handler, err := setupTestDB()
	if err != nil {
		b.Fatal(err)
	}

	// 自动迁移基准测试模型
	if err := handler.AutoMigrate(&BenchmarkUser{}); err != nil {
		b.Fatal(err)
	}

	// 插入测试数据
	users := make([]BenchmarkUser, 1000)
	for i := 0; i < 1000; i++ {
		users[i] = BenchmarkUser{
			Username:   "user" + string(rune(i)),
			Email:      "user" + string(rune(i)) + "@test.com",
			Age:        20 + (i % 50),
			BusinessID: int64(1 + (i % 10)),
			ShopID:     int64(100 + (i % 20)),
			Status:     1 + (i % 3),
			Tags:       "tag1,tag2,tag3",
		}
	}

	// 批量插入
	if err := handler.DB().CreateInBatches(users, 100).Error; err != nil {
		b.Fatal(err)
	}

	cleanup := func() {
		handler.Close()
	}

	return handler, cleanup
}

// BenchmarkSimpleQuery 简单查询基准测试
func BenchmarkSimpleQuery(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	param := NewSimpleQueryParam(BusinessIDQuery, 1)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var users []BenchmarkUser
			handler.Query(param).Find(&users)
		}
	})
}

// BenchmarkPageQuery 分页查询基准测试
func BenchmarkPageQuery(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	param := NewPageQueryParam(
		BusinessIDQuery,
		[]interface{}{1},
		10, // limit
		0,  // offset
		"created_at DESC",
	)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var users []BenchmarkUser
			handler.Query(param).Find(&users)
		}
	})
}

// BenchmarkQueryBuilder 查询构建器基准测试
func BenchmarkQueryBuilder(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			param := NewQueryBuilder().
				WithBusinessId(1).
				WithShopId(101).
				WhereIn("status", []interface{}{1, 2}).
				WithOrder("created_at", "DESC").
				WithPagination(10, 0).
				Build()

			var users []BenchmarkUser
			handler.Query(param).Find(&users)
		}
	})
}

// BenchmarkAdvancedQuery 高级查询基准测试
func BenchmarkAdvancedQuery(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	option := &FindOptionCommon{
		BusinessId:  1,
		Limit:       10,
		Offset:      0,
		By:          "created_at",
		Order:       "DESC",
		TablePrefix: "",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			param := NewAdvancedQueryParam(option)
			param.AddFilter(&BaseInfoFilter{
				DBField:    "status",
				Values:     []interface{}{1, 2},
				ExactMatch: true,
			})

			var users []BenchmarkUser
			handler.Query(param).Find(&users)
		}
	})
}

// BenchmarkTransaction 事务基准测试
func BenchmarkTransaction(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx := handler.Begin()

		user := BenchmarkUser{
			Username:   "bench_user",
			Email:      "bench@test.com",
			Age:        25,
			BusinessID: 1,
			ShopID:     101,
			Status:     1,
		}

		tx.DB().Create(&user)
		tx.Commit()
	}
}

// BenchmarkMultipleFilters 多过滤器基准测试
func BenchmarkMultipleFilters(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			param := NewQueryBuilder().
				WithBusinessId(1).
				WhereIn("status", []interface{}{1, 2}).
				WhereIn("age", []interface{}{25, 30, 35}).
				WhereLike("username", []interface{}{"user"}, false).
				WithOrder("created_at", "DESC").
				WithPagination(20, 0).
				Build()

			var users []BenchmarkUser
			handler.Query(param).Find(&users)
		}
	})
}

// BenchmarkQueryBuilderCreation 查询构建器创建基准测试
func BenchmarkQueryBuilderCreation(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NewQueryBuilder().
				WithBusinessId(1).
				WithShopId(101).
				WhereIn("status", []interface{}{1, 2}).
				WithOrder("created_at", "DESC").
				WithPagination(10, 0).
				Build()
		}
	})
}

// BenchmarkFilterCreation 过滤器创建基准测试
func BenchmarkFilterCreation(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NewInFilter("status", []interface{}{1, 2, 3})
			_ = NewLikeFilter("name", []interface{}{"test"}, true)
		}
	})
}

// BenchmarkComplexQuery 复杂查询基准测试
func BenchmarkComplexQuery(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	startTime := time.Now().AddDate(0, 0, -30).Format("2006-01-02 15:04:05")
	endTime := time.Now().Format("2006-01-02 15:04:05")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			param := NewQueryBuilder().
				WithBusinessId(1).
				WhereIn("status", []interface{}{1, 2}).
				WhereTimeRange("created_at", startTime, endTime).
				WhereFindInSet("tags", []string{"tag1", "tag2"}).
				WithGroupBy("business_id").
				WithOrder("created_at", "DESC").
				WithPagination(50, 0).
				Build()

			var users []BenchmarkUser
			handler.Query(param).Find(&users)
		}
	})
}

// BenchmarkBatchInsert 批量插入基准测试
func BenchmarkBatchInsert(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		users := make([]BenchmarkUser, 100)
		for j := 0; j < 100; j++ {
			users[j] = BenchmarkUser{
				Username:   "batch_user",
				Email:      "batch@test.com",
				Age:        25,
				BusinessID: 1,
				ShopID:     101,
				Status:     1,
			}
		}
		handler.DB().CreateInBatches(users, 50)
	}
}

// BenchmarkDirectGormQuery 直接 GORM 查询基准测试（对比）
func BenchmarkDirectGormQuery(b *testing.B) {
	handler, cleanup := setupBenchmarkDB(b)
	defer cleanup()

	db := handler.DB()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var users []BenchmarkUser
			db.Where(BusinessIDQuery, 1).Limit(10).Find(&users)
		}
	})
}
