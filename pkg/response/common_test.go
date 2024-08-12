package response

import (
	"testing"

	"github.com/spf13/cast"
)

func TestAllResponse(t *testing.T) {
	t.Run("TestGetBusinessCodeMsg", TestGetBusinessCodeMsg)
	t.Run("TestSetBusinessCode", TestSetBusinessCode)
}

func TestGetBusinessCodeMsg(t *testing.T) {
	expectedMsg := "成功"
	actualMsg := GetBusinessCodeMsg(SUCCESS)

	if expectedMsg != actualMsg {
		t.Errorf("Expected message: %s, but got: %s", expectedMsg, actualMsg)
	}
}

func TestSetBusinessCode(t *testing.T) {
	newBusinessCode := 999
	newMessage := "自定义消息"

	SetBusinessCode(BusinessCode(newBusinessCode), newMessage)

	actualMessage := GetBusinessCodeMsg(BusinessCode(newBusinessCode))

	if actualMessage != newMessage {
		t.Errorf("SetBusinessCode did not set the message correctly. Expected: %s, Actual: %s", newMessage, actualMessage)
	}
}

func TestGetErrorMsgAndBenchmarkBulk(t *testing.T) {
	businessCodes := []BusinessCode{SUCCESS, FAIL, ServerError}
	expectedMsgs := map[BusinessCode]string{
		SUCCESS:     "成功",
		FAIL:        "失败",
		ServerError: "服务器错误",
	}
	for _, BusinessCode := range businessCodes {
		expectedMsg := expectedMsgs[BusinessCode]

		t.Run("TestGetErrorMsg_"+cast.ToString(BusinessCode), func(t *testing.T) {
			actualMsg := GetBusinessCodeMsg(BusinessCode)

			if expectedMsg != actualMsg {
				t.Errorf("Expected message for %s: %s, but got: %s", cast.ToString(BusinessCode), expectedMsg, actualMsg)
			}
		})
	}
}
