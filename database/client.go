/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 13:55:35
 * @FilePath: \go-core\database\client.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package database

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/kamalyes/go-core/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Gorm 初始化数据库并产生数据库全局变量
// @return: *gorm.DB
func Gorm() *gorm.DB {
	switch global.CONFIG.Server.DataDriver {
	case "mysql":
		return GormMySQL()
	case "postgre":
		return GormPostgreSQL()
	case "sqlite":
		return GormSQLite()
	default:
		return GormMySQL()
	}
}

// GormMySQL 初始化Mysql数据库
// @return: *gorm.DB
func GormMySQL() *gorm.DB {
	config := global.CONFIG.MySQL
	if config.Host == "" {
		return nil
	}
	// 构建需要转义的参数值
	host := url.QueryEscape(config.Host)
	user := url.QueryEscape(config.Username)
	password := url.QueryEscape(config.Password)
	dbname := url.QueryEscape(config.Dbname)
	port := url.QueryEscape(config.Port)
	configString := url.QueryEscape(config.Config)
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?" + configString
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(config.LogLevel)); err != nil {
		global.LOG.Error("MySQL启动异常", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		if config.ConnMaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
		}
		if config.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
		}
		return db
	}
}

// gormConfig 根据配置决定是否开启日志
// @param: mod bool
// @return: *gorm.Config
func gormConfig(logLevel string) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}}
	switch logLevel {
	case "silent", "Silent":
		config.Logger = logger.Default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = logger.Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = logger.Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = logger.Default.LogMode(logger.Info)
	default:
		config.Logger = logger.Default.LogMode(logger.Error)
	}
	return config
}

// GormPostgreSQL 初始化PostgreSQL数据库
// @return: *gorm.DB
func GormPostgreSQL() *gorm.DB {
	config := global.CONFIG.PostgreSQL
	if config.Host == "" {
		return nil
	}
	// 构建需要转义的参数值
	host := url.QueryEscape(config.Host)
	user := url.QueryEscape(config.Username)
	password := url.QueryEscape(config.Password)
	dbname := url.QueryEscape(config.Dbname)
	port := url.QueryEscape(config.Port)
	configString := url.QueryEscape(config.Config)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s", host, user, password, dbname, port, configString)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // 禁用隐式 prepared statement
	}), gormConfig(config.LogLevel))
	if err != nil {
		global.LOG.Error("PostgreSQL启动异常", zap.Any("err", err))
		os.Exit(0)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		if config.ConnMaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
		}
		if config.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
		}
		return db
	}
}

// GormSQLite 连接sqlite 数据库
func GormSQLite() *gorm.DB {
	config := global.CONFIG.SQLite
	db, err := gorm.Open(sqlite.Open(config.DbPath), gormConfig(config.LogLevel))
	if err != nil {
		global.LOG.Error("SQLite数据库连接失败：", zap.Any("err", err))
		os.Exit(0)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		if config.ConnMaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
		}
		if config.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
		}
		return db
	}
}
