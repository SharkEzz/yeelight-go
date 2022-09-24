package device_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SharkEzz/yeelight-go/pkg/device"
	"github.com/SharkEzz/yeelight-go/pkg/device/command"
	"github.com/SharkEzz/yeelight-go/pkg/device/response"
)

// Mock the client to test the device
type FakeClient struct {
	isConnected bool
	response    chan *response.Response
}

func (c *FakeClient) Connect() error {
	c.isConnected = true
	c.response = make(chan *response.Response)
	return nil
}

func (c *FakeClient) Disconnect() error {
	c.isConnected = false
	close(c.response)
	return nil
}

func (c *FakeClient) SendCommand(cmd *command.Command) error {
	if !c.isConnected {
		return fmt.Errorf("client is not connected")
	}

	switch cmd.Method {
	case "set_bright":
	case "toggle":
	case "set_color":
		c.response <- &response.Response{
			ID:     cmd.ID,
			Result: []string{"ok"},
		}
	case "fake_notification":
		c.response <- &response.Response{
			ID: cmd.ID,
			Params: map[string]any{
				"power": "on",
			},
			ResponseType: "notification",
		}
	case "fake_error":
		c.response <- &response.Response{
			ID:    cmd.ID,
			Error: map[string]any{"code": -1, "message": "error"},
		}
	default:
		panic("unknown method")
	}

	return nil
}

func (c *FakeClient) GetResponse() <-chan *response.Response {
	return c.response
}

func (c *FakeClient) IsConnected() bool {
	return c.isConnected
}

func TestConnect(t *testing.T) {
	device := device.NewDevice(&FakeClient{})

	if err := device.Connect(); err != nil {
		t.Errorf("failed to connect: %v", err)
	}
	defer device.Disconnect()

	if !device.IsConnected() {
		t.Error("device is not connected")
	}
}

func TestSetBright(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		t.Log("timeout")
		t.FailNow()
	})
	light := device.NewDevice(&FakeClient{})

	if err := light.SetBright(100); err == nil {
		t.Error("device should not be connected")
	}

	if err := light.Connect(); err != nil {
		t.Error("failed to connect", err)
	}
	defer light.Disconnect()

	go func() {
		res := light.OnResponse()
		if res == nil || len(res.Result) != 1 || res.Result[0] != "ok" {
			t.Error("invalid response:", res)
		}
	}()

	if err := light.SetBright(255); err != nil {
		t.Error("there should be no error", err)
	}
}

func TestNotification(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		t.Log("timeout")
		t.FailNow()
	})
	light := device.NewDevice(&FakeClient{})

	if err := light.Connect(); err != nil {
		t.Error("failed to connect", err)
	}
	defer light.Disconnect()

	go func() {
		res := light.OnResponse()
		if res == nil || res.ResponseType != "notification" || res.Params["power"] != "on" {
			t.Error("invalid notification:", res)
		}
	}()

	err := light.SendCommand(&command.Command{
		ID:     1,
		Method: "fake_notification",
	})
	if err != nil {
		t.Error("there should be no error", err)
	}
}

func TestError(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		t.Log("timeout")
		t.FailNow()
	})
	light := device.NewDevice(&FakeClient{})

	if err := light.Connect(); err != nil {
		t.Error("failed to connect", err)
	}
	defer light.Disconnect()

	go func() {
		res := light.OnResponse()
		if res == nil || res.Error == nil || res.Error["code"] != -1 || res.Error["message"] != "error" {
			t.Error("invalid error:", res)
		}
	}()

	err := light.SendCommand(&command.Command{
		ID:     1,
		Method: "fake_error",
	})
	if err != nil {
		t.Error("there should be no error", err)
	}
}

func TestToggle(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		t.Log("timeout")
		t.FailNow()
	})
	light := device.NewDevice(&FakeClient{})

	if err := light.Connect(); err != nil {
		t.Error("failed to connect", err)
	}
	defer light.Disconnect()

	go func() {
		res := light.OnResponse()
		if res == nil || len(res.Result) != 1 || res.Result[0] != "ok" {
			t.Error("invalid response:", res)
		}
	}()

	err := light.Toggle()
	if err != nil {
		t.Error("there should be no error", err)
	}
}

func TestBright(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		t.Log("timeout")
		t.FailNow()
	})
	light := device.NewDevice(&FakeClient{})

	if err := light.Connect(); err != nil {
		t.Error("failed to connect", err)
	}
	defer light.Disconnect()

	go func() {
		res := light.OnResponse()
		if res == nil || len(res.Result) != 1 || res.Result[0] != "ok" {
			t.Error("invalid response:", res)
		}
	}()

	err := light.SetBright(255)
	if err != nil {
		t.Error("there should be no error", err)
	}
}

func TestColor(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		t.Log("timeout")
		t.FailNow()
	})
	light := device.NewDevice(&FakeClient{})

	if err := light.Connect(); err != nil {
		t.Error("failed to connect", err)
	}
	defer light.Disconnect()

	go func() {
		res := light.OnResponse()
		if res == nil || len(res.Result) != 1 || res.Result[0] != "ok" {
			t.Error("invalid response:", res)
		}
	}()

	err := light.SetColor(0xFFFFFF)
	if err != nil {
		t.Error("there should be no error", err)
	}
}
