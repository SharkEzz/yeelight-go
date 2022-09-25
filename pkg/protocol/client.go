package protocol

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/SharkEzz/yeelight-go/pkg/device/command"
	"github.com/SharkEzz/yeelight-go/pkg/device/response"
)

type YeelightClient struct {
	ip          string
	port        uint16
	isConnected bool
	conn        net.Conn
	Response    chan *response.Response
	lastCommand *command.Command
}

func NewYeelightClient(ip string, port uint16) *YeelightClient {
	client := &YeelightClient{
		ip:   ip,
		port: port,
	}

	return client
}

func (c *YeelightClient) Connect() error {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", c.ip, c.port))
	if err != nil {
		return err
	}

	c.isConnected = true
	c.conn = conn
	c.Response = make(chan *response.Response)

	go c.readResponse()

	return nil
}

func (c *YeelightClient) Disconnect() error {
	if !c.isConnected {
		return fmt.Errorf("client is not connected")
	}

	c.isConnected = false
	c.conn.Close()
	close(c.Response)

	return nil
}

func (c *YeelightClient) IsConnected() bool {
	return c.isConnected
}

// loop while connected, read response and send it to the response channel
func (c *YeelightClient) readResponse() {
	rd := bufio.NewReader(c.conn)

	for c.isConnected {
		data, err := rd.ReadString('\n')
		if err != nil {
			log.Println("error while reading response:", err)
		}

		response, err := response.ParseResponse(data)
		if err != nil {
			log.Println("error while parsing response", err)
		}

		c.Response <- response
	}
}

func (c *YeelightClient) GetResponse() <-chan *response.Response {
	return c.Response
}

func (c *YeelightClient) SendCommand(cmd *command.Command) error {
	if !c.isConnected {
		return fmt.Errorf("client is not connected")
	}

	cmdJson, err := cmd.GenerateJSON()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(cmdJson)
	if err != nil {
		return err
	}

	c.lastCommand = cmd

	return nil
}

func (c *YeelightClient) GetIP() string {
	return c.ip
}
