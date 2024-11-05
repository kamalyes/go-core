/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-03 20:15:09
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-03 20:22:33
 * @FilePath: \go-core\tests\response_echox_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/go-core/pkg/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenEchoResponse(t *testing.T) {
	e := echo.New()

	e.GET("/test", func(c echo.Context) error {
		respOption := &response.ResponseOption{
			Data:     "test data",
			Code:     200,
			Message:  "success",
			HttpCode: http.StatusOK,
		}
		return response.GenEchoResponse(c, respOption)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"data":"test data"`)
	assert.Contains(t, rec.Body.String(), `"code":200`)
	assert.Contains(t, rec.Body.String(), `"message":"success"`)
}
