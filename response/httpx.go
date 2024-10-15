package response

import (
	"encoding/json"
	"net/http"
)

// GenNetHttpResponse 生成 JSON 格式的响应
func GenNetHttpResponse(w http.ResponseWriter, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.merge()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(respOption.HttpCode))

	// 创建一个map来存储不包含HttpStatusCode和Language的字段
	cleanedResp := map[string]interface{}{
		"data":    respOption.Data,
		"code":    respOption.Code,
		"message": respOption.Message,
	}

	json.NewEncoder(w).Encode(cleanedResp)
}

// GenNetHttp400xResponse 生成 HTTP 400x 错误响应
func GenNetHttp400xResponse(w http.ResponseWriter, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = BadRequest
	respOption.HttpCode = StatusBadRequest
	GenNetHttpResponse(w, respOption)
}

// GenNetHttp500xResponse 生成 HTTP 500 错误响应
func GenNetHttp500xResponse(w http.ResponseWriter, respOption *ResponseOption) {
	if respOption == nil {
		respOption = &ResponseOption{}
	}
	respOption.Code = Fail
	respOption.HttpCode = StatusInternalServerError
	GenNetHttpResponse(w, respOption)
}
