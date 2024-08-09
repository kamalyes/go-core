/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-08 17:15:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 10:07:26
 * @FilePath: \go-core\i18n\i18n_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package i18n

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestConvertStringToTimestamp(t *testing.T) {
	expectedTimestamp := int64(1628424042)
	dateString := "2021-08-08 12:03:42"
	layout := "2006-01-02 15:05:05"
	timeZone := "UTC"

	timestamp, err := ConvertStringToTimestamp(dateString, layout, timeZone)

	assert.NoError(t, err)
	assert.Equal(t, expectedTimestamp, timestamp, "Timestamps should match")
}

func TestGetMsgWithMap(t *testing.T) {
	Init(language.Chinese)
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
	SetLang(cast.ToString(language.English))
	content = GetMsgWithMap(key, maps)
	assertEqual("XianMing AllName=Alice Age=30", content)
}

func TestGetMsgByKey(t *testing.T) {
	Init(language.Chinese)
	key := "test_key"
	content := GetMsgByKey(key)
	assert.NotEmpty(t, content)

}
