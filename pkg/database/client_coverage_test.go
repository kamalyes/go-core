/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 12:00:00
 * @FilePath: \go-core\pkg\database\client_coverage_test.go
 * @Description: client.go的覆盖率测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/kamalyes/go-config/pkg/database"
	"github.com/stretchr/testify/assert"
)

// TestBuildDSNAllTypes 测试buildDSN函数的所有数据库类型
func TestBuildDSNAllTypes(t *testing.T) {
	config := database.DBConfig{
		Host:     "localhost",
		Username: "testuser",
		Password: "testpass",
		Dbname:   "testdb",
		Port:     "3306",
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
		DbPath:   "/path/to/test.db",
	}

	// 测试MySQL DSN
	mysqlDSN := buildDSN(config, database.DBTypeMySQL)
	assert.Contains(t, mysqlDSN, "testuser")
	assert.Contains(t, mysqlDSN, "testpass")
	assert.Contains(t, mysqlDSN, "localhost")
	assert.Contains(t, mysqlDSN, "3306")
	assert.Contains(t, mysqlDSN, "testdb")
	assert.Contains(t, mysqlDSN, "charset")

	// 测试PostgreSQL DSN
	postgresDSN := buildDSN(config, database.DBTypePostgres)
	assert.Contains(t, postgresDSN, "host=localhost")
	assert.Contains(t, postgresDSN, "user=testuser")
	assert.Contains(t, postgresDSN, "password=testpass")
	assert.Contains(t, postgresDSN, "dbname=testdb")
	assert.Contains(t, postgresDSN, "port=3306")

	// 测试SQLite DSN
	sqliteDSN := buildDSN(config, database.DBTypeSQLite)
	assert.Equal(t, "/path/to/test.db", sqliteDSN)

	// 测试未知数据库类型
	unknownDSN := buildDSN(config, "unknown")
	assert.Equal(t, "", unknownDSN)
}

// TestGormConfigLogLevels 测试gormConfig函数的所有日志级别
func TestGormConfigLogLevels(t *testing.T) {
	// 测试silent级别
	config := gormConfig("silent")
	assert.NotNil(t, config)
	assert.True(t, config.DisableForeignKeyConstraintWhenMigrating)
	assert.NotNil(t, config.NamingStrategy)

	// 测试Silent级别（大写）
	config = gormConfig("Silent")
	assert.NotNil(t, config)

	// 测试error级别
	config = gormConfig("error")
	assert.NotNil(t, config)

	// 测试Error级别（大写）
	config = gormConfig("Error")
	assert.NotNil(t, config)

	// 测试warn级别
	config = gormConfig("warn")
	assert.NotNil(t, config)

	// 测试Warn级别（大写）
	config = gormConfig("Warn")
	assert.NotNil(t, config)

	// 测试info级别
	config = gormConfig("info")
	assert.NotNil(t, config)

	// 测试Info级别（大写）
	config = gormConfig("Info")
	assert.NotNil(t, config)

	// 测试默认级别
	config = gormConfig("unknown")
	assert.NotNil(t, config)

	// 测试空字符串
	config = gormConfig("")
	assert.NotNil(t, config)
}

// TestBuildDSNSpecialCharacters 测试buildDSN处理特殊字符
func TestBuildDSNSpecialCharacters(t *testing.T) {
	config := database.DBConfig{
		Host:     "test@host",
		Username: "user&name",
		Password: "pass=word",
		Dbname:   "db/name",
		Port:     "5432",
		Config:   "sslmode=disable&TimeZone=Asia/Shanghai",
		DbPath:   "/path/to/test with spaces.db",
	}

	// 测试MySQL DSN转义
	mysqlDSN := buildDSN(config, database.DBTypeMySQL)
	assert.Contains(t, mysqlDSN, "user%26name") // &被转义
	assert.Contains(t, mysqlDSN, "pass%3Dword") // =被转义
	assert.Contains(t, mysqlDSN, "test%40host") // @被转义
	assert.Contains(t, mysqlDSN, "db%2Fname")   // /被转义

	// 测试PostgreSQL DSN转义
	postgresDSN := buildDSN(config, database.DBTypePostgres)
	assert.Contains(t, postgresDSN, "test%40host") // @被转义
	assert.Contains(t, postgresDSN, "user%26name") // &被转义

	// 测试SQLite DSN（不需要转义）
	sqliteDSN := buildDSN(config, database.DBTypeSQLite)
	assert.Equal(t, "/path/to/test with spaces.db", sqliteDSN)
}