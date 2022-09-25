package protocol

import (
	"github.com/SharkEzz/yeelight-go/pkg/device/command"
	"github.com/SharkEzz/yeelight-go/pkg/device/response"
)

// Client is the interface that wraps the basic methods of a Yeelight client.
//
// The reception of responses is done through the channel returned by the GetResponse method.
type Client interface {
	Connect() error
	Disconnect() error
	SendCommand(cmd *command.Command) error
	GetResponse() <-chan *response.Response
	IsConnected() bool
	GetIP() string
}
