package response

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenNetHttpResponse(t *testing.T) {
	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 创建一个响应选项
	respOption := &ResponseOption{
		Data:     "test data",
		Code:     200,
		Message:  "success",
		HttpCode: http.StatusOK,
	}

	// 调用 GenNetHttpResponse
	GenNetHttpResponse(w, respOption)

	// 检查响应状态码
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// 检查响应内容
	expectedBody := `{"data":"test data","code":200,"message":"success"}`
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.JSONEq(t, expectedBody, string(body))
}

func TestGenNetHttp400xResponse(t *testing.T) {
	w := httptest.NewRecorder()

	respOption := &ResponseOption{
		Message: "bad request",
	}

	// 调用 GenNetHttp400xResponse
	GenNetHttp400xResponse(w, respOption)

	// 检查响应状态码
	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// 检查响应内容
	expectedBody := `{"data":null,"code":400,"message":"bad request"}`
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.JSONEq(t, expectedBody, string(body))
}

func TestGenNetHttp500xResponse(t *testing.T) {
	w := httptest.NewRecorder()

	respOption := &ResponseOption{
		Message: "internal server error",
	}

	// 调用 GenNetHttp500xResponse
	GenNetHttp500xResponse(w, respOption)

	// 检查响应状态码
	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	// 检查响应内容
	expectedBody := `{"data":null,"code":500,"message":"internal server error"}`
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.JSONEq(t, expectedBody, string(body))
}
