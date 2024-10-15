/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-08 17:15:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 17:20:08
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
	t.Run("TestGetMsgWithMap", TestGetMsgWithMap)
	t.Run("TestGetMsgByKey", TestGetMsgByKey)
	t.Run("TestIsValidLanguage", TestIsValidLanguage)
	t.Run("TestUseI18nMiddleware", TestUseI18nMiddleware)
	t.Run("TestSetCurrentLang", TestSetCurrentLang)
	t.Run("TestAddLanguage", TestAddLanguage)
	t.Run("TestSetCurrentLanguage", TestSetCurrentLanguage)
	t.Run("TestMessageRetrieval", TestMessageRetrieval)
	t.Run("TestRemoveLanguage", TestRemoveLanguage)

}

func TestGetMsgWithMap(t *testing.T) {
	LoadMessageFileFS(Chinese)
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
	LoadMessageFileFS(Language("zh-CN"))
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
		invalid := IsValidLanguage(Language("zh-TW"))
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

func TestAddLanguage(t *testing.T) {
	customLanguage := Language("fr")

	// 测试添加新语言
	AddLanguage(Language(customLanguage))
	assert.Equal(t, IsValidLanguage(customLanguage), true)
}

func TestSetCurrentLanguage(t *testing.T) {
	var err error
	customLanguage := Language("bg")

	// 测试将当前语言设置为自定义语言
	LoadMessageFileFS(customLanguage)
	err = SetCurrentLang(customLanguage)
	assert.Nil(t, err)
	assert.Equal(t, customLanguage, currentLang)
}

func TestMessageRetrieval(t *testing.T) {
	customLanguage := Language("bg")
	AddLanguage(Language(customLanguage))
	LoadMessageFileFS(customLanguage)
	key := "test_key"
	maps := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}

	// 测试使用和不使用映射来检索消息
	content := GetMsgByKey(key)
	assert.Equal(t, "XianMing AllName=<no value> Възраст=<no value>", content)

	content = GetMsgWithMap(key, maps)
	assert.Equal(t, "XianMing AllName=Alice Възраст=30", content)
}

func TestRemoveLanguage(t *testing.T) {
	// 初始化测试条件
	SetEarsLang(Language("bg"))
	SetCurrentLang(Language("en"))

	// 测试移除当前语言
	err := RemoveLanguage(Language(GetCurrentLang()))
	assert.NotNil(t, err)
	// assert.Equal(t, fmt.Sprintf("cannot remove the currently selected language %s", GetCurrentLang()), err.Error())

	// 测试移除兜底语言
	err = RemoveLanguage(Language(GetEarsLang()))
	assert.NotNil(t, err)
	// assert.Equal(t, fmt.Sprintf("cannot remove the ears language %s", GetEarsLang()), err.Error())

	// 测试移除存在的支持语言
	toRemoveLang := Language("zh-CN")
	err = RemoveLanguage(toRemoveLang)
	assert.Nil(t, err)

	// 测试移除不存在的语言
	langNotPresent := Language("de")
	err = RemoveLanguage(langNotPresent)
	assert.Nil(t, err)
}
