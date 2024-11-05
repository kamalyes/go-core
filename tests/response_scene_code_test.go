package tests

import (
	"testing"

	"github.com/kamalyes/go-core/pkg/response"
	"github.com/spf13/cast"
)

func TestAllResponse(t *testing.T) {
	t.Run("TestGetSceneCodeMsg", TestGetSceneCodeMsg)
	t.Run("TestSetSceneCode", TestSetSceneCode)
}

func TestGetSceneCodeMsg(t *testing.T) {
	expectedMsg := "Success"
	actualMsg := response.GetSceneCodeMsg(response.Success)

	if expectedMsg != actualMsg {
		t.Errorf("Expected message: %s, but got: %s", expectedMsg, actualMsg)
	}
}

func TestSetSceneCode(t *testing.T) {
	newSceneCode := 999
	newMessage := "自定义消息"

	response.SetSceneCode(response.SceneCode(newSceneCode), newMessage)

	actualMessage := response.GetSceneCodeMsg(response.SceneCode(newSceneCode))

	if actualMessage != newMessage {
		t.Errorf("SetSceneCode did not set the message correctly. Expected: %s, Actual: %s", newMessage, actualMessage)
	}
}

func TestGetErrorMsgAndBenchmarkBulk(t *testing.T) {
	businessCodes := []response.SceneCode{response.Success, response.Fail, response.ServerError}
	expectedMsgs := map[response.SceneCode]string{
		response.Success:     "Success",
		response.Fail:        "Fail",
		response.ServerError: "Internal Server Error",
	}
	for _, SceneCode := range businessCodes {
		expectedMsg := expectedMsgs[SceneCode]

		t.Run("TestGetErrorMsg_"+cast.ToString(SceneCode), func(t *testing.T) {
			actualMsg := response.GetSceneCodeMsg(SceneCode)

			if expectedMsg != actualMsg {
				t.Errorf("Expected message for %s: %s, but got: %s", cast.ToString(SceneCode), expectedMsg, actualMsg)
			}
		})
	}
}
