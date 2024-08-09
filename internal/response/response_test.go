/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-09 13:55:26
 * @FilePath: \go-core\internal\response\response_test.go
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
	"github.com/kamalyes/go-core/internal/code"
)

func TestResult(t *testing.T) {
	// 设置默认的 httpCode
	SetHTTPCode(http.StatusNotFound)
	testCases := []struct {
		name     string
		code     code.StatusCode
		content  interface{}
		message  string
		expected int
	}{
		{"Test 1", http.StatusOK, "Data", "Success", http.StatusOK},
		{"Test 2", http.StatusOK, "Data", "Success", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			Result(ctx, tc.code, tc.content, tc.message)

			if w.Code != tc.expected {
				t.Errorf("Expected %s status code: %d, got: %d", tc.name, tc.expected, w.Code)
			}
		})
	}
}

func TestOk(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	Ok(ctx)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code: %d, got: %d", http.StatusOK, w.Code)
	}
}
