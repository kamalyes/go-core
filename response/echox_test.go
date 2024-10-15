package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenEchoResponse(t *testing.T) {
	e := echo.New()

	e.GET("/test", func(c echo.Context) error {
		respOption := &ResponseOption{
			Data:     "test data",
			Code:     200,
			Message:  "success",
			HttpCode: http.StatusOK,
		}
		return GenEchoResponse(c, respOption)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"data":"test data"`)
	assert.Contains(t, rec.Body.String(), `"code":200`)
	assert.Contains(t, rec.Body.String(), `"message":"success"`)
}
