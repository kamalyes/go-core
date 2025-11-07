/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 10:28:19
 * @FilePath: \go-core\pkg\database\page.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package database

import (
	"bytes"
	"errors"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/pkg/global"
	"github.com/kamalyes/go-toolbox/pkg/convert"
)

const (

	/** ------- and 条件 ------  */
	/** 小于 < */
	lt = "lt:"

	/** 大于 > */
	gt = "gt:"

	/** 小于 <= */
	lte = "lte:"

	/** 大于 >= */
	gte = "gte:"

	/** 默认是等于 */
	eq = "eq:"

	/** 模糊查询 */
	lk = "lk:"

	/** ------- or 条件 ------  */

	/** 小于 */
	orlt = "orlt:"

	/** 大于  */
	orgt = "orgt:"

	/** 小于 */
	orlte = "orlte:"

	/** 大于  */
	orgte = "orgte:"

	/** 默认是等于 */
	oreq = "oreq:"

	/** 模糊查询 */
	orlk = "orlk:"

	/** ------- 排序 ------  */

	/** 降序 */
	pd = ":pd:"

	/** 升序 */
	pa = ":pa:"
)

// PageBean 全局分页对象
type PageBean struct {

	/** 当前页  */
	Page int `json:"page"`

	/** 当前页的行数 */
	PageSize int `json:"pageSize"`

	/** 总记录数 */
	Total int64 `json:"total"`

	/** 每行的数据 */
	Rows interface{} `json:"rows"`
}

type PageInfo struct {

	/** 当前页 */
	Current int

	/** 每页显示的最大行数 */
	RowCount int

	/** 表名 仅限于指定表名去查询 */
	TableName string

	/** 查询 and 条件参数 */
	AndParams map[string]interface{}

	/** 查询 or 条件参数 */
	OrParams map[string]interface{}

	/** 排序 */
	OrderStr string
}

// parseBasicParams 解析基础分页参数
func parseBasicParams(key, value string, pageInfo *PageInfo) bool {
	switch key {
	case "current":
		current, err := strconv.Atoi(value)
		if err != nil {
			current = 1
		}
		if current < 1 {
			current = 1
		}
		pageInfo.Current = current
		return true
	case "rowCount":
		rowCount, err := strconv.Atoi(value)
		if err != nil {
			rowCount = 10
		}
		if rowCount < 1 {
			rowCount = 10
		} else if rowCount > 100 {
			rowCount = 100
		}
		pageInfo.RowCount = rowCount
		return true
	case "orderStr":
		pageInfo.OrderStr = value
		return true
	case "tableName":
		pageInfo.TableName = value
		return true
	}
	return false
}

// parseAndCondition 解析AND条件参数
func parseAndCondition(key, value string, andParams map[string]interface{}) bool {
	key = CamelToCase(key)
	
	switch {
	case strings.Index(value, lt) == 0:
		value = strings.Replace(value, lt, "", 1)
		if value != "" {
			andParams[key+" < ?"] = value
		}
	case strings.Index(value, lte) == 0:
		value = strings.Replace(value, lte, "", 1)
		if value != "" {
			andParams[key+" <= ?"] = value
		}
	case strings.Index(value, gt) == 0:
		value = strings.Replace(value, gt, "", 1)
		if value != "" {
			andParams[key+" > ?"] = value
		}
	case strings.Index(value, gte) == 0:
		value = strings.Replace(value, gte, "", 1)
		if value != "" {
			andParams[key+" >= ?"] = value
		}
	case strings.Index(value, lk) == 0:
		value = strings.Replace(value, lk, "", 1)
		if value != "" {
			andParams[key+" LIKE ?"] = value + "%"
		}
	case strings.Index(value, eq) == 0:
		value = strings.Replace(value, eq, "", 1)
		if value != "" {
			andParams[key+" = ?"] = value
		}
	default:
		// 默认等于条件
		if value != "" {
			andParams[key+" = ?"] = value
		}
	}
	return true
}

// parseOrCondition 解析OR条件参数
func parseOrCondition(key, value string, orParams map[string]interface{}) bool {
	key = CamelToCase(key)
	
	switch {
	case strings.Index(value, orlt) == 0:
		value = strings.Replace(value, orlt, "", 1)
		if value != "" {
			orParams[key+" < ?"] = value
		}
	case strings.Index(value, orlte) == 0:
		value = strings.Replace(value, orlte, "", 1)
		if value != "" {
			orParams[key+" <= ?"] = value
		}
	case strings.Index(value, orgte) == 0:
		value = strings.Replace(value, orgte, "", 1)
		if value != "" {
			orParams[key+" >= ?"] = value
		}
	case strings.Index(value, orgt) == 0:
		value = strings.Replace(value, orgt, "", 1)
		if value != "" {
			orParams[key+" > ?"] = value
		}
	case strings.Index(value, orlk) == 0:
		value = strings.Replace(value, orlk, "", 1)
		if value != "" {
			orParams[key+" LIKE ?"] = value + "%"
		}
	case strings.Index(value, oreq) == 0:
		value = strings.Replace(value, oreq, "", 1)
		if value != "" {
			orParams[key+" = ?"] = value
		}
	default:
		return false
	}
	return true
}

// processOrderString 处理排序字符串
func processOrderString(orderStr string) string {
	if orderStr == "" {
		return ""
	}
	v := CamelToCase(orderStr)
	v = strings.ReplaceAll(v, pd, " desc,")
	v = strings.ReplaceAll(v, pa, " asc,")
	v = strings.TrimSuffix(v, ",")
	return v
}

// PageParam 获取url查询参数
func PageParam(c *gin.Context) *PageInfo {
	s := c.Request.URL.RawQuery
	paramStr, err := url.QueryUnescape(s)
	if err != nil {
		log.Println("url参数decode异常：" + err.Error())
		return nil
	}
	
	pageInfo := PageInfo{}
	andParams := make(map[string]interface{})
	orParams := make(map[string]interface{})
	
	paramArr := strings.Split(paramStr, "&")
	for _, v := range paramArr {
		ky := strings.Split(v, "=")
		if len(ky) != 2 {
			continue
		}
		
		key := ky[0]
		value := ky[1]
		
		// 处理基础分页参数
		if parseBasicParams(key, value, &pageInfo) {
			continue
		}
		
		// 跳过时间戳参数
		if key == "_t" || key == "_time" || key == "_timestamp" {
			continue
		}
		
		// 先尝试处理OR条件
		if parseOrCondition(key, value, orParams) {
			continue
		}
		
		// 处理AND条件（包括默认等于条件）
		parseAndCondition(key, value, andParams)
	}
	
	// 处理排序字符串
	pageInfo.OrderStr = processOrderString(pageInfo.OrderStr)
	pageInfo.AndParams = andParams
	pageInfo.OrParams = orParams
	return &pageInfo
}

func CamelToCase(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// Buffer 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case string:
		b.append(val)
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case []byte:
		_, _ = b.Write(val)
	case rune:
		_, _ = b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	_, _ = b.WriteString(s)
	return b
}

// CheckPageRows 获取页数和行数
func CheckPageRows(currentStr, rowCountStr string) (current, rowCount int) {
	current, err := strconv.Atoi(currentStr)
	if err != nil {
		current = 1
	}
	if current < 1 {
		current = 1
	}
	rowCount, err = strconv.Atoi(rowCountStr)
	if err != nil {
		rowCount = 10
	}
	if rowCount < 1 {
		rowCount = 10
	} else if rowCount > 500 {
		rowCount = 500
	}
	return current, rowCount
}

// FindPage 分页查询 v-空对象指针
func FindPage(v interface{}, rows interface{}, pageInfo *PageInfo) (pageBean *PageBean, err error) {
	if pageInfo == nil {
		return nil, errors.New("入参pageInfo不能为空指针")
	}
	pageBean = &PageBean{Page: pageInfo.Current, PageSize: pageInfo.RowCount}
	var total int64
	db := global.DB.Model(v)
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() == reflect.String {
		db = global.DB.Table(convert.MustString(v))
	}
	andCons := pageInfo.AndParams
	orCons := pageInfo.OrParams
	orderStr := pageInfo.OrderStr
	if len(andCons) > 0 {
		for k, v := range andCons {
			db = db.Where(k, v)
		}
	}
	if len(orCons) > 0 {
		for k, v := range orCons {
			db = db.Or(k, v)
		}
	}
	db.Count(&total)
	if len(orderStr) > 0 {
		err = db.Limit(pageBean.PageSize).Offset((pageBean.Page - 1) * pageBean.PageSize).Order(orderStr).Find(rows).Error
	} else {
		err = db.Limit(pageBean.PageSize).Offset((pageBean.Page - 1) * pageBean.PageSize).Find(rows).Error
	}
	pageBean.Rows = rows
	pageBean.Total = total
	return
}
