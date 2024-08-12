/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-10 23:41:12
 * @FilePath: \go-core\pkg\response\ginx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSendResponse(t *testing.T) {
	// 创建一个模拟的 Gin 上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建要测试的 ResponseOption
	responseOption := ResponseOption{}

	// 调用发送响应函数
	responseOption.Sub(c)

	// 检查响应状态码是否正确
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	// 在此添加其他相应字段的测试，比如JSON响应的内容等
}

func TestResponse(t *testing.T) {
	testCases := []struct {
		name     string
		httpCode int
		code     BusinessCode
		data     interface{}
		message  string
		expected int
	}{
		{"Test 1", 0, 0, "", "", 200},
		{"Test 2", 200, http.StatusOK, "Data", "Success", http.StatusOK},
		{"Test 3", 404, http.StatusNotFound, "Data", "StatusNotFound", http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			option := ResponseOption{
				Code:     tc.code,
				HttpCode: tc.httpCode,
				Data:     tc.data,
				Message:  tc.message,
			}
			GenResponse(ctx, &option)

			if w.Code != tc.expected {
				t.Errorf("Expected %s status code: %d, got: %d", tc.name, tc.expected, w.Code)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	// 创建初始 ResponseOption
	initialRespOption := &ResponseOption{}

	// 创建要合并的 ResponseOption
	mergeRespOption := &ResponseOption{
		Code:     300,
		Data:     "Merged Data",
		Message:  "Merged Message",
		HttpCode: 500,
	}

	// 执行 merge 操作
	initialRespOption.merge(mergeRespOption)

	// 检查合并后的字段是否符合预期
	if initialRespOption.Code != mergeRespOption.Code {
		t.Errorf("Expected merged code to be %d, but got %d", mergeRespOption.Code, initialRespOption.Code)
	}

	if initialRespOption.Data != mergeRespOption.Data {
		t.Errorf("Expected merged data to be %s, but got %s", mergeRespOption.Data, initialRespOption.Data)
	}

	if initialRespOption.Message != mergeRespOption.Message {
		t.Errorf("Expected message to remain %s, but got %s", mergeRespOption.Message, initialRespOption.Message)
	}

	if initialRespOption.HttpCode != mergeRespOption.HttpCode {
		t.Errorf("Expected HTTP code to remain %d, but got %d", mergeRespOption.HttpCode, initialRespOption.HttpCode)
	}
}
