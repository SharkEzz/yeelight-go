package protocol

import (
	"github.com/SharkEzz/yeelight-go/pkg/device/command"
	"github.com/SharkEzz/yeelight-go/pkg/device/response"
)

type Client interface {
	Connect() error
	Disconnect() error
	SendCommand(cmd *command.Command) error
	GetResponse() <-chan *response.Response
	IsConnected() bool
}
