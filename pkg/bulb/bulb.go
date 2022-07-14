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
	Name          string
	IP            string
	isConnected   bool
	conn          net.Conn
	reader        *bufio.Reader
	lastCommandId int
}

// Create a new Bulb instance with no socket opened.
func NewBulb(name, ip string) *Bulb {
	return &Bulb{
		name,
		ip,
		false,
		nil,
		nil,
		0,
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
	return nil
}

// Close the connection to the bulb
func (b *Bulb) Disconnect() error {
	return b.conn.Close()
}

// Send a command to the bulb.
//
// Prefer to use it with predefined command function like `SetBright`, `SetHSV`, etc... to avoid issues.
func (l *Bulb) SendCommand(cmd *command.Command) (*response.Response, error) {
	if !l.isConnected {
		return nil, fmt.Errorf("bulb is not connected")
	}

	cmdJson, err := cmd.GenerateJSON()
	if err != nil {
		return nil, err
	}

	_, err = l.conn.Write(cmdJson)
	if err != nil {
		return nil, err
	}
	l.lastCommandId = cmd.ID

	data, err := l.reader.ReadString('\n')
	if err != nil {
		return nil, nil
	}

	response, err := response.ParseResponse(data)
	if err != nil {
		return nil, err
	}

	return response, nil
}
