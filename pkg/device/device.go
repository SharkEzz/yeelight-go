package device

import (
	"bufio"
	"fmt"
	"net"

	"github.com/SharkEzz/yeelight-go/pkg/device/command"
	"github.com/SharkEzz/yeelight-go/pkg/device/response"
)

// The bulb struct represent a connection to a Yeelight light bulb
type Device struct {
	Name            string
	IP              string
	isConnected     bool
	ResponseChannel chan *response.Response
	LastCommandId   int
	conn            net.Conn
	reader          *bufio.Reader
}

// Create a new Bulb instance with no socket opened.
func NewBulb(name, ip string) *Device {
	return &Device{
		Name:        name,
		IP:          ip,
		isConnected: false,
	}
}

// Open a socket to the bulb with the adress and port specified when instantiating a new Bulb
func (d *Device) Connect() error {
	if d.isConnected {
		return fmt.Errorf("connection is already established")
	}

	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:55443", d.IP))
	if err != nil {
		return err
	}

	d.isConnected = true
	d.conn = conn
	d.reader = bufio.NewReader(conn)

	d.ResponseChannel = make(chan *response.Response, 256)

	go d.readLoop()

	return nil
}

// Close the connection to the bulb
func (d *Device) Disconnect() error {
	if !d.isConnected {
		return fmt.Errorf("connection is not established")
	}

	d.isConnected = false
	d.reader = nil
	close(d.ResponseChannel)
	return d.conn.Close()
}

func (d *Device) IsConnected() bool {
	return d.isConnected
}

// Send a command to the bulb.
//
// Prefer to use it with predefined command function like `SetBright`, `SetHSV`, etc... to avoid issues.
func (d *Device) SendCommand(cmd *command.Command) error {
	if !d.isConnected {
		return fmt.Errorf("bulb is not connected")
	}

	cmdJson, err := cmd.GenerateJSON()
	if err != nil {
		return err
	}

	_, err = d.conn.Write(cmdJson)
	if err != nil {
		return err
	}
	d.LastCommandId = cmd.ID

	return nil
}

// Loop while the bulb is connected, and read the responses from the bulb before sending
// them to the ResponseChannel.
func (d *Device) readLoop() {
	for d.isConnected {
		data, err := d.reader.ReadString('\n')
		if err != nil {
			return
		}

		response, err := response.ParseResponse(data)
		if err != nil {
			return
		}

		d.ResponseChannel <- response
	}
}

/* COMMANDS */

func (d *Device) GetProp(props []any) (int, error) {
	cmd := command.NewCommand(command.GET_PROP, props)

	err := d.SendCommand(cmd)
	if err != nil {
		return 0, err
	}

	return cmd.ID, nil
}

// Set the bulb brightness.
//
// The brightness is a value between 0 and 100.
func (d *Device) SetBright(brightness uint8) error {
	if brightness < 1 {
		brightness = 1
	} else if brightness > 100 {
		brightness = 100
	}

	err := d.SendCommand(command.NewCommand(command.SET_BRIGHT, []any{brightness, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Set the bulb color to the specified RGB values.
//
// The RGB values are values between 0 and 255.
func (d *Device) SetRGB(r, g, b uint32) error {
	if r > 255 || g > 255 || b > 255 {
		return fmt.Errorf("RGB values must be between 0 and 255")
	}

	rgb := (r << 16) | (g << 8) | b

	err := d.SendCommand(command.NewCommand(command.SET_RGB, []any{rgb, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Set the bulb color to the specified Hue and saturation values.
func (d *Device) SetHSV(h uint16, s uint8) error {
	if s > 100 {
		s = 100
	}

	err := d.SendCommand(command.NewCommand(command.SET_HSV, []any{h, s, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Set the bulb power state.
func (d *Device) SetPower(power bool) error {
	var powerStr string
	if power {
		powerStr = "on"
	} else {
		powerStr = "off"
	}

	err := d.SendCommand(command.NewCommand(command.SET_POWER, []any{powerStr, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Toggle the bulb power state.
func (d *Device) Toggle() error {
	err := d.SendCommand(command.NewCommand(command.TOGGLE, []any{}))
	if err != nil {
		return err
	}

	return nil
}
