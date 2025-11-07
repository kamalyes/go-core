package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSetupTestDB 测试setupTestDB函数
func TestSetupTestDB(t *testing.T) {
	db, handler, err := setupTestDB()
	assert.NoError(t, err)
	defer handler.Close()

	// 打印所有表
	tables, err := db.Migrator().GetTables()
	assert.NoError(t, err)
	t.Logf("Tables: %v", tables)

	// 检查表是否存在
	hasUsersTable := db.Migrator().HasTable(&TestUser{})
	hasProductsTable := db.Migrator().HasTable(&TestProduct{})

	t.Logf("Has users table: %v", hasUsersTable)
	t.Logf("Has products table: %v", hasProductsTable)

	assert.True(t, hasUsersTable, "test_users table should exist")
	assert.True(t, hasProductsTable, "test_products table should exist")
}