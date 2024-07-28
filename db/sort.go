/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:13:14
 * @FilePath: \go-core\db\crud.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package db

import "strings"

const (

	// DESC  降序
	DESC = "desc"

	// ASC  升序
	ASC = "asc"
)

type Sort struct {

	/** 排序字段 */
	OrderKey string `json:"orderKey"`

	/** 排序方式 */
	Type string `json:"type"`
}

// ResolveSortList 解析 SortList
func ResolveSortList(sl []*Sort) string {
	if len(sl) < 1 {
		return ""
	}
	if len(sl) == 1 {
		return sl[0].OrderKey + " " + sl[0].Type
	}
	sb := strings.Builder{}
	for i := range sl {
		sb.WriteString(sl[i].OrderKey)
		sb.WriteString(" ")
		sb.WriteString(sl[i].Type)
		if i == len(sl)-1 {
			break
		}
		sb.WriteString(",")
	}
	return sb.String()
}
