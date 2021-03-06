package command_test

import (
	"strings"
	"testing"

	"github.com/SharkEzz/yeelight-go/pkg/device/command"
)

func TestCommand_GenerateJSON(t *testing.T) {
	cmd1 := command.NewCommand(command.SET_RGB, []any{0xFF, 0xFF, 0xFF}, 1)
	data, err := cmd1.GenerateJSON()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.HasSuffix(string(data), "\r\n") {
		t.Errorf("Expected data to have \\r\\n suffix, got %s", string(data))
	}

	if string(data) != `{"id":1,"method":"set_rgb","params":[255,255,255]}`+"\r\n" {
		t.Errorf("Expected data to be %s, got %s", `{"id":1,"method":"set_rgb","params":[255,255,255]}`, string(data))
	}
}

func TestNewCommand(t *testing.T) {
	cmd1 := command.NewCommand(command.SET_RGB, []any{0xFF, 0xFF, 0xFF})
	cmd2 := command.NewCommand(command.SET_POWER, []any{"off"})

	if cmd1.ID == cmd2.ID {
		t.Errorf("Expected command ID to be different, got %d", cmd1.ID)
	}

	cmd3 := command.NewCommand(command.SET_POWER, []any{"off"}, 1)
	if cmd3.ID != 1 {
		t.Errorf("Expected command ID to be 1, got %d", cmd3.ID)
	}
}
