package bulb

import (
	"bufio"
	"fmt"
	"net"

	"github.com/SharkEzz/yeelight-go/pkg/bulb/command"
	"github.com/SharkEzz/yeelight-go/pkg/bulb/response"
)

// The bulb struct represent a connection to a Yeelight light bulb
type Bulb struct {
	Name            string
	IP              string
	isConnected     bool
	ResponseChannel chan *response.Response
	LastCommandId   int
	conn            net.Conn
	reader          *bufio.Reader
}

// Create a new Bulb instance with no socket opened.
func NewBulb(name, ip string) *Bulb {
	return &Bulb{
		Name:        name,
		IP:          ip,
		isConnected: false,
	}
}

// Open a socket to the bulb with the adress and port specified when instantiating a new Bulb
func (b *Bulb) Connect() error {
	if b.isConnected {
		return fmt.Errorf("connection is already established")
	}

	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:55443", b.IP))
	if err != nil {
		return err
	}

	b.isConnected = true
	b.conn = conn
	b.reader = bufio.NewReader(conn)

	b.ResponseChannel = make(chan *response.Response, 256)

	go b.readLoop()

	return nil
}

// Close the connection to the bulb
func (b *Bulb) Disconnect() error {
	if !b.isConnected {
		return fmt.Errorf("connection is not established")
	}

	b.isConnected = false
	b.reader = nil
	close(b.ResponseChannel)
	return b.conn.Close()
}

func (b *Bulb) IsConnected() bool {
	return b.isConnected
}

// Send a command to the bulb.
//
// Prefer to use it with predefined command function like `SetBright`, `SetHSV`, etc... to avoid issues.
func (l *Bulb) SendCommand(cmd *command.Command) error {
	if !l.isConnected {
		return fmt.Errorf("bulb is not connected")
	}

	cmdJson, err := cmd.GenerateJSON()
	if err != nil {
		return err
	}

	_, err = l.conn.Write(cmdJson)
	if err != nil {
		return err
	}
	l.LastCommandId = cmd.ID

	return nil
}

// Loop while the bulb is connected, and read the responses from the bulb before sending
// them to the ResponseChannel.
func (b *Bulb) readLoop() {
	for b.isConnected {
		data, err := b.reader.ReadString('\n')
		if err != nil {
			return
		}

		response, err := response.ParseResponse(data)
		if err != nil {
			return
		}

		b.ResponseChannel <- response
	}
}

/* COMMANDS */

func (b *Bulb) GetProp(props []any) (int, error) {
	cmd := command.NewCommand(command.GET_PROP, props)

	err := b.SendCommand(cmd)
	if err != nil {
		return 0, err
	}

	return cmd.ID, nil
}

// Set the bulb brightness.
//
// The brightness is a value between 0 and 100.
func (b *Bulb) SetBright(brightness uint8) error {
	if brightness < 1 {
		brightness = 1
	} else if brightness > 100 {
		brightness = 100
	}

	err := b.SendCommand(command.NewCommand(command.SET_BRIGHT, []any{brightness, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Set the bulb color to the specified RGB values.
//
// The RGB values are values between 0 and 255.
func (bu *Bulb) SetRGB(r, g, b uint32) error {
	if r > 255 || g > 255 || b > 255 {
		return fmt.Errorf("RGB values must be between 0 and 255")
	}

	rgb := (r << 16) | (g << 8) | b

	err := bu.SendCommand(command.NewCommand(command.SET_RGB, []any{rgb, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Set the bulb color to the specified Hue and saturation values.
func (b *Bulb) SetHSV(h uint16, s uint8) error {
	if s > 100 {
		s = 100
	}

	err := b.SendCommand(command.NewCommand(command.SET_HSV, []any{h, s, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Set the bulb power state.
func (b *Bulb) SetPower(power bool) error {
	var powerStr string
	if power {
		powerStr = "on"
	} else {
		powerStr = "off"
	}

	err := b.SendCommand(command.NewCommand(command.SET_POWER, []any{powerStr, "smooth", 500}))
	if err != nil {
		return err
	}

	return nil
}

// Toggle the bulb power state.
func (b *Bulb) Toggle() error {
	err := b.SendCommand(command.NewCommand(command.TOGGLE, []any{}))
	if err != nil {
		return err
	}

	return nil
}
