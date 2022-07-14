package command

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Initialize the random generator.
func init() {
	rand.Seed(time.Now().Unix())
}

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
func NewCommand(method Method, params []any, id ...int) *Command {
	var commandId int
	if len(id) == 0 {
		commandId = rand.Intn(9999)
	} else {
		commandId = id[0]
	}

	return &Command{
		commandId,
		method,
		params,
	}
}

// Ensure that an effect string is one of "smooth" or "sudden".
func validateEffect(effect string) (string, error) {
	e := strings.ToLower(effect)
	if e != "sudden" && e != "smooth" {
		return "", fmt.Errorf("invalid effect string provided")
	}

	return e, nil
}

// Create a new command for setting the bulb brightness.
func SetBright(brightness uint8, effect string, duration uint) *Command {
	effect, err := validateEffect(effect)
	if err != nil {
		effect = "smooth"
	}

	if brightness < 1 {
		brightness = 1
	} else if brightness > 100 {
		brightness = 100
	}

	return NewCommand(SET_BRIGHT, []any{brightness, effect, duration})
}
