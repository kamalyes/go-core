package response

import (
	"testing"

	"github.com/spf13/cast"
)

func TestAllResponse(t *testing.T) {
	t.Run("TestGetSceneCodeMsg", TestGetSceneCodeMsg)
	t.Run("TestSetSceneCode", TestSetSceneCode)
}

func TestGetSceneCodeMsg(t *testing.T) {
	expectedMsg := "Success"
	actualMsg := GetSceneCodeMsg(Success)

	if expectedMsg != actualMsg {
		t.Errorf("Expected message: %s, but got: %s", expectedMsg, actualMsg)
	}
}

func TestSetSceneCode(t *testing.T) {
	newSceneCode := 999
	newMessage := "自定义消息"

	SetSceneCode(SceneCode(newSceneCode), newMessage)

	actualMessage := GetSceneCodeMsg(SceneCode(newSceneCode))

	if actualMessage != newMessage {
		t.Errorf("SetSceneCode did not set the message correctly. Expected: %s, Actual: %s", newMessage, actualMessage)
	}
}

func TestGetErrorMsgAndBenchmarkBulk(t *testing.T) {
	businessCodes := []SceneCode{Success, Fail, ServerError}
	expectedMsgs := map[SceneCode]string{
		Success:     "Success",
		Fail:        "Fail",
		ServerError: "Internal Server Error",
	}
	for _, SceneCode := range businessCodes {
		expectedMsg := expectedMsgs[SceneCode]

		t.Run("TestGetErrorMsg_"+cast.ToString(SceneCode), func(t *testing.T) {
			actualMsg := GetSceneCodeMsg(SceneCode)

			if expectedMsg != actualMsg {
				t.Errorf("Expected message for %s: %s, but got: %s", cast.ToString(SceneCode), expectedMsg, actualMsg)
			}
		})
	}
}
