/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-13 16:15:55
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
	"github.com/kamalyes/go-core/pkg/httpx"
	"github.com/stretchr/testify/assert"
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
		name           string
		httpStatusCode httpx.StatusCode
		code           SceneCode
		data           interface{}
		message        string
		expected       int
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
				Code:       tc.code,
				StatusCode: tc.httpStatusCode,
				Data:       tc.data,
				Message:    tc.message,
			}
			GenResponse(ctx, &option)

			if w.Code != tc.expected {
				t.Errorf("Expected %s status code: %d, got: %d", tc.name, tc.expected, w.Code)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	// 初始化一个Response
	initRespOption := &ResponseOption{}

	// 执行 merge 操作
	initRespOption.merge()

	// 检查合并后的字段是否符合预期
	if initRespOption.Code != SceneCode(httpx.StatusOK) {
		t.Errorf("Expected initRespOption merged code to be %d", initRespOption.Code)
	}

	if initRespOption.Data != nil {
		t.Errorf("Expected initRespOption merged data to be %s", initRespOption.Data)
	}

	if initRespOption.Message != GetSceneCodeMsg(SceneCode(httpx.StatusOK)) {
		t.Errorf("Expected initRespOption message to remain %s", initRespOption.Message)
	}

	if initRespOption.StatusCode != Success {
		t.Errorf("Expected initRespOption HTTP code to remain %d", initRespOption.StatusCode)
	}

	// 定义新的自定义模型
	newCode := SceneCode(36666)
	newHttpStatusCode := httpx.StatusInternalServerError
	mergeRespOption := &ResponseOption{
		Code:       newCode,
		StatusCode: newHttpStatusCode,
	}
	mergeRespOption.merge()
	if mergeRespOption.Code != newCode || mergeRespOption.StatusCode != newHttpStatusCode {
		t.Errorf("Expected mergeRespOption HTTP code to remain %d", mergeRespOption.StatusCode)
	}

	if mergeRespOption.Message != httpx.GetStatusCodeText(newHttpStatusCode) {
		t.Errorf("Expected mergeRespOption message to remain %s", mergeRespOption.Message)
	}

	newMessage := "1235678"
	mergeRespOption.Message = newMessage
	mergeRespOption.merge()
	if mergeRespOption.Message != newMessage {
		t.Errorf("Expected mergeRespOption message to remain %s", mergeRespOption.Message)
	}

}

func TestGen400xResponse(t *testing.T) {
	router := gin.Default()

	router.GET("/testGen400xResponse", func(c *gin.Context) {
		Gen400xResponse(c, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testGen400xResponse", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 校验响应体内容
	expectedJSON := "{\"data\":null,\"code\":400,\"message\":\"Bad Request\"}"
	assert.JSONEq(t, expectedJSON, w.Body.String())
}

func TestGen500xResponse(t *testing.T) {
	router := gin.Default()

	router.GET("/testGen500xResponse", func(c *gin.Context) {
		Gen500xResponse(c, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testGen500xResponse", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 校验响应体内容
	expectedJSON := "{\"data\":null,\"code\":500,\"message\":\"Fail\"}"
	assert.JSONEq(t, expectedJSON, w.Body.String())
}
