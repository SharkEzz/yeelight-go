package response

import (
	"encoding/json"
)

// The response struct represent a bulb response.
type Response struct {
	ID           int
	Result       []string
	Error        map[string]any
	responseType string
}

// Attemp to guess the response type (error, command or notification) from the response content.
func guessResponseType(responseStruct *Response) string {
	if responseStruct.ID != 0 && len(responseStruct.Error) != 0 {
		return "error"
	} else if responseStruct.ID != 0 && len(responseStruct.Result) != 0 {
		return "command"
	}

	return "notification"
}

// Parse a raw JSON string response and return a Response instance.
func ParseResponse(response string) (*Response, error) {
	responseStruct := &Response{}

	err := json.Unmarshal([]byte(response), responseStruct)
	if err != nil {
		return nil, err
	}

	responseStruct.responseType = guessResponseType(responseStruct)

	return responseStruct, nil
}

// Return the response type.
func (r *Response) ResponseType() string {
	return r.responseType
}
