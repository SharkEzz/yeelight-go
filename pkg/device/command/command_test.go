package command_test

import (
	"strings"
	"testing"

	"github.com/SharkEzz/yeelight-go/pkg/device/command"
)

func TestCommand_GenerateJSON(t *testing.T) {
	cmd1 := command.NewCommand(command.SET_RGB, []any{0xFF, 0xFF, 0xFF})
	data, err := cmd1.GenerateJSON()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if strings.HasSuffix(string(data), `{"method":"set_rgb","params":[255,255,255]}`+"\r\n") {
		t.Errorf("Expected data to be %s, got %s", `{"method":"set_rgb","params":[255,255,255]}`, string(data))
	}
}
