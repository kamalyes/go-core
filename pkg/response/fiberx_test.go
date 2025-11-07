/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-03 20:15:09
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-03 20:15:09
 * @FilePath: \go-core\response_fiberx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package response

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/assert"
)

func TestGenFiberResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		respOption := &ResponseOption{
			Data:     "test data",
			Code:     200,
			Message:  "success",
			HttpCode: http.StatusOK,
		}
		return GenFiberResponse(c, respOption)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 读取响应体
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"data":"test data"`)
	assert.Contains(t, string(body), `"code":200`)
	assert.Contains(t, string(body), `"message":"success"`)
}
