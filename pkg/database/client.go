/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-03 20:02:19
 * @FilePath: \go-core\pkg\database\client.go
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

	"github.com/kamalyes/go-config/pkg/database"
	"github.com/kamalyes/go-core/pkg/global"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {
	switch global.CONFIG.Server.DataDriver {
	case database.DBTypeMySQL:
		return GormMySQL()
	case database.DBTypePostgres:
		return GormPostgreSQL()
	case database.DBTypeSQLite:
		return GormSQLite()
	default:
		return GormMySQL() // 默认使用 MySQL
	}
}

// GormMySQL 初始化MySQL数据库
func GormMySQL() *gorm.DB {
	config := global.CONFIG.MySQL
	dbConfig, err := database.NewDBConfig(database.DBTypeMySQL, config)
	if err != nil {
		global.LOGGER.WithError(err).ErrorMsg("MySQL config error")
		return nil
	}
	return initDB(*dbConfig, database.DBTypeMySQL, func(dsn string) (*gorm.DB, error) {
		return gorm.Open(mysql.New(mysql.Config{DSN: dsn}), gormConfig(dbConfig.LogLevel))
	})
}

// GormPostgreSQL 初始化PostgreSQL数据库
func GormPostgreSQL() *gorm.DB {
	config := global.CONFIG.PostgreSQL
	dbConfig, err := database.NewDBConfig(database.DBTypePostgres, config)
	if err != nil {
		global.LOGGER.WithError(err).ErrorMsg("PostgreSQL config error")
		return nil
	}
	return initDB(*dbConfig, database.DBTypePostgres, func(dsn string) (*gorm.DB, error) {
		return gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), gormConfig(dbConfig.LogLevel))
	})
}

// GormSQLite 连接SQLite数据库
func GormSQLite() *gorm.DB {
	config := global.CONFIG.SQLite
	dbConfig, err := database.NewDBConfig(database.DBTypeSQLite, config)
	if err != nil {
		global.LOGGER.WithError(err).ErrorMsg("SQLite config error")
		return nil
	}
	return initDB(*dbConfig, database.DBTypeSQLite, func(dsn string) (*gorm.DB, error) {
		return gorm.Open(sqlite.Open(dbConfig.DbPath), gormConfig(dbConfig.LogLevel))
	})
}

// initDB 初始化数据库连接
func initDB(config database.DBConfig, dbType string, openFunc func(string) (*gorm.DB, error)) *gorm.DB {
	if config.Host == "" {
		global.LOGGER.Error("Database host is empty")
		return nil
	}

	dsn := buildDSN(config, dbType)
	db, err := openFunc(dsn)
	if err != nil {
		global.LOGGER.ErrorKV(fmt.Sprintf("%s database startup error", dbType), "err", err)
		os.Exit(0)
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	if config.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
	}
	if config.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
	}

	return db
}

// buildDSN 构建数据库连接字符串
func buildDSN(config database.DBConfig, dbType string) string {
	host := url.QueryEscape(config.Host)
	user := url.QueryEscape(config.Username)
	password := url.QueryEscape(config.Password)
	dbname := url.QueryEscape(config.Dbname)
	port := url.QueryEscape(config.Port)
	configString := url.QueryEscape(config.Config)

	var dsn string
	switch dbType {
	case database.DBTypeMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, dbname, configString)
	case database.DBTypePostgres:
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s", host, user, password, dbname, port, configString)
	case database.DBTypeSQLite:
		dsn = config.DbPath
	}
	return dsn
}

// gormConfig 根据配置决定是否开启日志
func gormConfig(logLevel string) *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

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
