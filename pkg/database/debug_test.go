package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTableCreation 测试表创建
func TestTableCreation(t *testing.T) {
	db, handler, err := setupTestDB()
	assert.NoError(t, err)
	defer handler.Close()

	// 检查表是否存在
	hasUsersTable := db.Migrator().HasTable(&TestUser{})
	hasProductsTable := db.Migrator().HasTable(&TestProduct{})

	t.Logf("Has users table: %v", hasUsersTable)
	t.Logf("Has products table: %v", hasProductsTable)

	assert.True(t, hasUsersTable, "test_users table should exist")
	assert.True(t, hasProductsTable, "test_products table should exist")

	// 测试插入数据
	user := TestUser{
		Username:   "test_user",
		Email:      "test@example.com",
		Age:        25,
		BusinessID: 1,
		ShopID:     100,
		Status:     1,
		Tags:       "test",
	}

	err = db.Create(&user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// 验证数据
	var count int64
	err = db.Model(&TestUser{}).Where("username = ?", "test_user").Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	t.Logf("User count: %d", count)
}