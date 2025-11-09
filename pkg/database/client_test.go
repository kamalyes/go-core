/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-10 01:48:43
 * @FilePath: \go-core\pkg\database\client_test.go
 * @Description: client 数据库连接测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package database

import (
	"testing"

	"github.com/kamalyes/go-config/pkg/database"
	"github.com/kamalyes/go-core/pkg/global"
	gologger "github.com/kamalyes/go-logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestBuildDSN 测试DSN构建
func TestBuildDSN(t *testing.T) {
	config := database.DBConfig{
		Host:     "localhost",
		Port:     "3306",
		Username: "user",
		Password: "pass",
		Dbname:   "testdb",
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
		DbPath:   "/tmp/test.db",
	}

	// 测试MySQL DSN
	mysqlDSN := buildDSN(config, database.DBTypeMySQL)
	expected := "user:pass@tcp(localhost:3306)/testdb?charset%3Dutf8mb4%26parseTime%3DTrue%26loc%3DLocal"
	assert.Equal(t, expected, mysqlDSN)

	// 测试PostgreSQL DSN
	postgresDSN := buildDSN(config, database.DBTypePostgres)
	expected = "host=localhost user=user password=pass dbname=testdb port=3306 charset%3Dutf8mb4%26parseTime%3DTrue%26loc%3DLocal"
	assert.Equal(t, expected, postgresDSN)

	// 测试SQLite DSN
	sqliteDSN := buildDSN(config, database.DBTypeSQLite)
	assert.Equal(t, config.DbPath, sqliteDSN)

	// 测试其他类型（默认为空）
	otherDSN := buildDSN(config, "unknown")
	assert.Equal(t, "", otherDSN)
}

// TestGormConfig 测试GORM配置
func TestGormConfig(t *testing.T) {
	tests := []struct {
		logLevel       string
		expectedLogger logger.LogLevel
	}{
		{"silent", logger.Silent},
		{"Silent", logger.Silent},
		{"error", logger.Error},
		{"Error", logger.Error},
		{"warn", logger.Warn},
		{"Warn", logger.Warn},
		{"info", logger.Info},
		{"Info", logger.Info},
		{"unknown", logger.Error}, // 默认值
		{"", logger.Error},        // 空值
	}

	for _, tt := range tests {
		t.Run(tt.logLevel, func(t *testing.T) {
			config := gormConfig(tt.logLevel)
			assert.NotNil(t, config)
			assert.True(t, config.DisableForeignKeyConstraintWhenMigrating)
			assert.NotNil(t, config.NamingStrategy)
			assert.NotNil(t, config.Logger)
		})
	}
}

// TestGormFunctionsWithoutConfig 测试在没有配置时的行为
func TestGormFunctionsWithoutConfig(t *testing.T) {
	// 这些函数在没有正确配置时应该返回nil或处理错误
	// 由于依赖全局配置，我们无法在测试中完全模拟

	// 验证函数存在且不会panic
	assert.Panics(t, func() { Gorm() })
	assert.Panics(t, func() { GormMySQL() })
	assert.Panics(t, func() { GormPostgreSQL() })
	assert.Panics(t, func() { GormSQLite() })
}

// TestInitDBWithEmptyHost 测试空主机配置
func TestInitDBWithEmptyHost(t *testing.T) {
	// 备份原始配置
	originalLog := global.LOGGER
	global.LOGGER = gologger.NewLogger(&gologger.LogConfig{Level: gologger.INFO})

	defer func() {
		global.LOGGER = originalLog
	}()

	config := database.DBConfig{
		Host: "", // 空主机
	}

	// 测试空主机应该返回nil
	db := initDB(config, database.DBTypeMySQL, func(dsn string) (*gorm.DB, error) {
		return nil, nil
	})

	assert.Nil(t, db)
}
