/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-08 17:15:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 17:28:54
 * @FilePath: \go-core\i18n\i18n.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package i18n

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	I18n              = new(i18n.Localizer)
	defaultDateFormat = "2006-01-02 15:05:05"
	earsLang          = Language("en") // 兜底语言
	currentLang       = earsLang       // 全局变量用来存储lang值
	langMutex         sync.Mutex
	//go:embed lang/*
	fs     embed.FS
	bundle *i18n.Bundle
)

type Language string

const (
	English Language = "en"    // 英文
	Chinese Language = "zh-CN" // 简体中文
)

var SupportedLanguages = map[Language]struct{}{
	English: {},
	Chinese: {},
}

// IsValidLanguage 检查给定的语言是否受支持
func IsValidLanguage(lang Language) bool {
	_, ok := SupportedLanguages[lang]
	return ok
}

// AddLanguage 添加支持的语言
func AddLanguage(lang Language) {
	// 检查语言若不存在则添加
	if _, ok := SupportedLanguages[lang]; !ok {
		SupportedLanguages[lang] = struct{}{}
	}
}

// RemoveLanguage 移除支持的语言
func RemoveLanguage(lang Language) error {
	// 如果要移除的语言是当前使用的语言，则返回错误
	if lang == currentLang {
		return fmt.Errorf("cannot remove the currently selected language %s", lang)
	}

	// 若要移除的语言是兜底语言
	if lang == earsLang {
		return fmt.Errorf("cannot remove the ears language %s", lang)
	}

	// 从支持的语言中删除指定语言
	if _, ok := SupportedLanguages[lang]; ok {
		delete(SupportedLanguages, lang)
		LoadMessageFileFS()
	}
	return nil
}

// SetCurrentLang 设置当前语言
func SetCurrentLang(lang Language) error {
	AddLanguage(lang)
	langMutex.Lock()
	defer langMutex.Unlock()
	currentLang = lang
	return nil
}

// GetCurrentLang 获取当前语言
func GetCurrentLang() string {
	return string(currentLang)
}

// SetCurrentLang 设置兜底语言
func SetEarsLang(lang Language) error {
	AddLanguage(lang)
	langMutex.Lock()
	defer langMutex.Unlock()
	earsLang = lang
	return nil
}

// GetEarsLang 获取兜底语言
func GetEarsLang() string {
	return string(earsLang)
}

func (lang Language) ToTextLanguage() language.Tag {
	return language.Make(string(lang))
}

// ConvertStringToTimestamp String时间类型转换为时间戳
func ConvertStringToTimestamp(dateString, layout string, timeZone string) (int64, error) {
	// 加载时区
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return 0, err
	}

	// 将字符串解析为时间
	t, err := time.ParseInLocation(layout, dateString, loc)
	if err != nil {
		return 0, err
	}
	// 转换时间为时间戳
	timestamp := t.Unix()
	return timestamp, nil
}

// GetMsgWithMap 获取变量值
func GetMsgWithMap(key string, maps map[string]interface{}) string {
	I18n := i18n.NewLocalizer(bundle, GetCurrentLang())
	var content string
	if maps == nil {
		content, _ = I18n.Localize(&i18n.LocalizeConfig{
			MessageID: key,
		})
	} else {
		content, _ = I18n.Localize(&i18n.LocalizeConfig{
			MessageID:    key,
			TemplateData: maps,
		})
	}
	content = strings.ReplaceAll(content, ": <no value>", "")
	if content == "" {
		return key
	} else {
		return content
	}
}

// GetMsgByKey 通过Key获取内容
func GetMsgByKey(key string) string {
	I18n := i18n.NewLocalizer(bundle, GetCurrentLang())
	content, err := I18n.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		// 处理错误，例如记录日志或返回默认消息
		fmt.Printf("GetMsgByKey Error %s", err)
	}
	return content
}

// GetTimeOffset 国际化时间戳偏移
func GetTimeOffset(timezone string, ts int64) (offset int) {
	var loc, _ = time.LoadLocation(timezone)
	_, offset = time.Unix(ts, 0).In(loc).Zone()
	return
}

// FormatWithLocation 国际化时间戳转换字符串
func FormatWithLocation(timezone string, ts int64) string {
	lt, _ := time.LoadLocation(timezone)
	str := time.Unix(ts, 0).In(lt).Format(defaultDateFormat)
	return str
}

// ParseWithLocation 国际化时间字符串转换时间戳
func ParseWithLocation(timezone string, timeStr string) int64 {
	l, _ := time.LoadLocation(timezone)
	lt, _ := time.ParseInLocation(defaultDateFormat, timeStr, l)
	return lt.Unix()
}

// UseI18n gin框架使用
func UseI18n() gin.HandlerFunc {
	return func(context *gin.Context) {
		acceptLang := context.GetHeader("Accept-Language")
		if _, ok := SupportedLanguages[Language(acceptLang)]; ok {
			currentLang = Language(acceptLang)
		} else {
			fmt.Printf("Request language (%s) not supported, current language (%s)", acceptLang, currentLang)
		}
		LoadMessageFileFS(Language(acceptLang))
	}
}

// LoadMessageFileFS 其它初始化i18n Bundle并加载支持的语言文件
func LoadMessageFileFS(lang ...Language) {
	if len(lang) > 0 {
		currentLang = Language(lang[0])
	}
	SetCurrentLang(currentLang)
	bundle = i18n.NewBundle(currentLang.ToTextLanguage())
	I18n = i18n.NewLocalizer(bundle, GetCurrentLang())
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	// 加载所有支持的语言文件
	for lang := range SupportedLanguages {
		langFile := fmt.Sprintf("%s.yaml", lang)
		_, err := bundle.LoadMessageFileFS(fs, "lang/"+langFile)
		if err != nil {
			fmt.Printf("failed loading translation for %s, err: %v", lang, err)
			os.Exit(2)
		}
	}
}
