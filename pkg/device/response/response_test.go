package response_test

import (
	"testing"

	"github.com/SharkEzz/yeelight-go/pkg/device/response"
)

func TestParseResponse(t *testing.T) {
	res := `{"id":1,"result":["ok"],"error":null}` + "\r\n"
	responseStruct, err := response.ParseResponse(res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if responseStruct.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", responseStruct.ID)
	}

	if responseStruct.ResponseType != "command" {
		t.Errorf("Expected response type to be command, got %s", responseStruct.ResponseType)
	}

	if len(responseStruct.Result) != 1 {
		t.Errorf("Expected result to be 1, got %d", len(responseStruct.Result))
	}

	if responseStruct.Result[0] != "ok" {
		t.Errorf("Expected result to be ok, got %s", responseStruct.Result[0])
	}

	res = `{"id":1,"result":["ok"],"error":{"code":1,"message":"error message"}}` + "\r\n"
	responseStruct, err = response.ParseResponse(res)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if responseStruct.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", responseStruct.ID)
	}

	if responseStruct.ResponseType != "error" {
		t.Errorf("Expected response type to be error, got %s", responseStruct.ResponseType)
	}

	if len(responseStruct.Error) != 2 {
		t.Errorf("Expected error to be 2, got %d", len(responseStruct.Error))
	}

	if responseStruct.Error["code"] != 1. {
		t.Errorf("Expected error code to be 1, got %f", responseStruct.Error["code"])
	}

	if responseStruct.Error["message"] != "error message" {
		t.Errorf("Expected error message to be error message, got %s", responseStruct.Error["message"])
	}

	res = `{"id":1,"method":"props","params":{"brightness":100,"color_temp":6500}}` + "\r\n"
	responseStruct, err = response.ParseResponse(res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if responseStruct.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", responseStruct.ID)
	}

	if responseStruct.ResponseType != "notification" {
		t.Errorf("Expected response type to be command, got %s", responseStruct.ResponseType)
	}

	if len(responseStruct.Params) != 2 {
		t.Errorf("Expected result to be 2, got %d", len(responseStruct.Result))
	}
}
