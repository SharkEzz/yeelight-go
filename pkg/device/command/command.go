package command

import (
	"encoding/json"
	"math/rand"
)

// The command struct represent a command to be send to the bulb.
type Command struct {
	ID     int    `json:"id"`
	Method Method `json:"method"`
	Params []any  `json:"params"`
}

// Generate a correctly formatted JSON string that can be send to the bulb.
func (c *Command) GenerateJSON() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	data = append(data, "\r\n"...)

	return data, nil
}

// Create a new command with the given method and params.
//
// Prefer to use predefined function to generate commands.
func NewCommand(method Method, params []any) *Command {
	commandId := rand.Intn(9999)

	return &Command{
		commandId,
		method,
		params,
	}
}
