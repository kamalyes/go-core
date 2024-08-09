/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-08 17:15:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 13:37:50
 * @FilePath: \go-core\i18n\i18n_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAllI18nFunctions(t *testing.T) {
	t.Run("TestConvertStringToTimestamp", TestConvertStringToTimestamp)
	t.Run("TestFormatWithLocation", TestFormatWithLocation)
	t.Run("TestGetMsgWithMap", TestGetMsgWithMap)
	t.Run("TestGetMsgByKey", TestGetMsgByKey)
	t.Run("TestIsValidLanguage", TestIsValidLanguage)
	t.Run("TestUseI18nMiddleware", TestUseI18nMiddleware)
	t.Run("TestSetCurrentLang", TestSetCurrentLang)
}

func TestConvertStringToTimestamp(t *testing.T) {
	expectedTimestamp := int64(1628424042)
	dateString := "2021-08-08 12:03:42"
	layout := "2006-01-02 15:05:05"
	timeZone := "UTC"

	timestamp, err := ConvertStringToTimestamp(dateString, layout, timeZone)

	assert.NoError(t, err)
	assert.Equal(t, expectedTimestamp, timestamp, "Timestamps should match")
}

func TestFormatWithLocation(t *testing.T) {
	expected := "2024-05-16 13:09:09"
	timestamp := int64(1715867289) // This timestamp represents "2024-05-16 21:58:09" in UTC
	formatted := FormatWithLocation("UTC", timestamp)
	assert.Equal(t, expected, formatted)
}

func TestGetMsgWithMap(t *testing.T) {
	Init(Chinese)
	key := "test_key"
	maps := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}

	assertEqual := func(expected, content string) {
		assert.Equal(t, expected, content, "expected error")
	}

	content := GetMsgWithMap(key, maps)
	assertEqual("小明全名=Alice,年龄=30", content)
	// 重新赋值
	SetCurrentLang(Language("en"))
	content = GetMsgWithMap(key, maps)
	assertEqual("XianMing AllName=Alice Age=30", content)
}

func TestGetMsgByKey(t *testing.T) {
	Init(Language("zh-CN"))
	key := "test_key"
	content := GetMsgByKey(key)
	assert.NotEmpty(t, content)
}

func TestIsValidLanguage(t *testing.T) {
	t.Run("Valid Language", func(t *testing.T) {
		valid := IsValidLanguage(English)
		assert.True(t, valid)
	})

	t.Run("Invalid Language", func(t *testing.T) {
		invalid := IsValidLanguage(Language("fr"))
		assert.False(t, invalid)
	})
}

func TestUseI18nMiddleware(t *testing.T) {
	router := gin.New()
	router.GET("/test", UseI18n(), func(c *gin.Context) {
		msg := GetMsgByKey("test_key")
		c.String(http.StatusOK, msg)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Language", "zh-CN")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "小明全名=<no value>,年龄=<no value>", w.Body.String())
}

func TestSetCurrentLang(t *testing.T) {
	SetCurrentLang("zh-CN")
	assert.Equal(t, Language("zh-CN"), currentLang)
}
