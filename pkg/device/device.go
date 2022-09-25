package device

import (
	"github.com/SharkEzz/yeelight-go/pkg/device/command"
	"github.com/SharkEzz/yeelight-go/pkg/device/response"
	"github.com/SharkEzz/yeelight-go/pkg/protocol"
)

type Device struct {
	client protocol.Client
}

func NewDevice(client protocol.Client) *Device {
	return &Device{
		client: client,
	}
}

func (d *Device) IsConnected() bool {
	return d.client.IsConnected()
}

func (d *Device) Connect() error {
	return d.client.Connect()
}

func (d *Device) Disconnect() error {
	return d.client.Disconnect()
}

func (d *Device) SendCommand(cmd *command.Command) error {
	return d.client.SendCommand(cmd)
}

func (d *Device) OnResponse() *response.Response {
	return <-d.client.GetResponse()
}

func (d *Device) GetIP() string {
	return d.client.GetIP()
}

// Commands part

func (d *Device) SetBright(bright uint8) (int, error) {
	cmd := command.NewCommand(command.SET_BRIGHT, []any{bright, "smooth", 500})

	return cmd.ID, d.client.SendCommand(cmd)
}

func (d *Device) SetColor(color uint32) (int, error) {
	cmd := command.NewCommand(command.SET_RGB, []any{color, "smooth", 500})

	return cmd.ID, d.client.SendCommand(cmd)
}

func (d *Device) Toggle() (int, error) {
	cmd := command.NewCommand(command.TOGGLE, []any{})

	return cmd.ID, d.client.SendCommand(cmd)
}

func (d *Device) SetPower(power bool) (int, error) {
	var powerStr bool
	if power {
		powerStr = true
	}

	cmd := command.NewCommand(command.SET_POWER, []any{powerStr, "smooth", 500})

	return cmd.ID, d.client.SendCommand(cmd)
}

func (d *Device) SetHSV(hue, sat uint16) (int, error) {
	cmd := command.NewCommand(command.SET_HSV, []any{hue, sat, "smooth", 500})

	return cmd.ID, d.client.SendCommand(cmd)
}

func (d *Device) GetProp(prop string) (int, error) {
	cmd := command.NewCommand(command.GET_PROP, []any{prop})

	return cmd.ID, d.client.SendCommand(cmd)
}
